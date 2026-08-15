"""RMVPE 人声音高估计（RVC MIT 实现精简：仅 torch.stft 推理路径）。"""
from __future__ import annotations

from typing import List

import numpy as np
import torch
import torch.nn as nn
import torch.nn.functional as F
from librosa.filters import mel as librosa_mel


class BiGRU(nn.Module):
    def __init__(self, input_features: int, hidden_features: int, num_layers: int):
        super().__init__()
        self.gru = nn.GRU(
            input_features,
            hidden_features,
            num_layers=num_layers,
            batch_first=True,
            bidirectional=True,
        )

    def forward(self, x: torch.Tensor) -> torch.Tensor:
        return self.gru(x)[0]


class ConvBlockRes(nn.Module):
    def __init__(self, in_channels: int, out_channels: int, momentum: float = 0.01):
        super().__init__()
        self.conv = nn.Sequential(
            nn.Conv2d(in_channels, out_channels, (3, 3), padding=1, bias=False),
            nn.BatchNorm2d(out_channels, momentum=momentum),
            nn.ReLU(),
            nn.Conv2d(out_channels, out_channels, (3, 3), padding=1, bias=False),
            nn.BatchNorm2d(out_channels, momentum=momentum),
            nn.ReLU(),
        )
        self.shortcut = (
            nn.Conv2d(in_channels, out_channels, (1, 1)) if in_channels != out_channels else None
        )

    def forward(self, x: torch.Tensor) -> torch.Tensor:
        if self.shortcut is None:
            return self.conv(x) + x
        return self.conv(x) + self.shortcut(x)


class ResEncoderBlock(nn.Module):
    def __init__(
        self,
        in_channels: int,
        out_channels: int,
        kernel_size,
        n_blocks: int = 1,
        momentum: float = 0.01,
    ):
        super().__init__()
        self.conv = nn.ModuleList()
        self.conv.append(ConvBlockRes(in_channels, out_channels, momentum))
        for _ in range(n_blocks - 1):
            self.conv.append(ConvBlockRes(out_channels, out_channels, momentum))
        self.kernel_size = kernel_size
        if kernel_size is not None:
            self.pool = nn.AvgPool2d(kernel_size=kernel_size)

    def forward(self, x: torch.Tensor):
        for conv in self.conv:
            x = conv(x)
        if self.kernel_size is not None:
            return x, self.pool(x)
        return x


class Encoder(nn.Module):
    def __init__(
        self,
        in_channels: int,
        in_size: int,
        n_encoders: int,
        kernel_size,
        n_blocks: int,
        out_channels: int = 16,
        momentum: float = 0.01,
    ):
        super().__init__()
        self.n_encoders = n_encoders
        self.bn = nn.BatchNorm2d(in_channels, momentum=momentum)
        self.layers = nn.ModuleList()
        self.latent_channels = []
        for i in range(self.n_encoders):
            self.layers.append(
                ResEncoderBlock(
                    in_channels, out_channels, kernel_size, n_blocks, momentum=momentum
                )
            )
            self.latent_channels.append([out_channels, in_size])
            in_channels = out_channels
            out_channels *= 2
            in_size //= 2
        self.out_size = in_size
        self.out_channel = out_channels

    def forward(self, x: torch.Tensor):
        concat_tensors: List[torch.Tensor] = []
        x = self.bn(x)
        for i in range(self.n_encoders):
            t, x = self.layers[i](x)
            concat_tensors.append(t)
        return x, concat_tensors


class Intermediate(nn.Module):
    def __init__(
        self,
        in_channels: int,
        out_channels: int,
        n_inters: int,
        n_blocks: int,
        momentum: float = 0.01,
    ):
        super().__init__()
        self.n_inters = n_inters
        self.layers = nn.ModuleList()
        self.layers.append(ResEncoderBlock(in_channels, out_channels, None, n_blocks, momentum))
        for _ in range(self.n_inters - 1):
            self.layers.append(
                ResEncoderBlock(out_channels, out_channels, None, n_blocks, momentum)
            )

    def forward(self, x: torch.Tensor) -> torch.Tensor:
        for i in range(self.n_inters):
            x = self.layers[i](x)
        return x


class ResDecoderBlock(nn.Module):
    def __init__(self, in_channels: int, out_channels: int, stride, n_blocks: int = 1, momentum: float = 0.01):
        super().__init__()
        out_pad = (0, 1) if stride == (1, 2) else (1, 1)
        self.conv1 = nn.Sequential(
            nn.ConvTranspose2d(
                in_channels, out_channels, kernel_size=(3, 3), stride=stride, padding=1, output_padding=out_pad, bias=False
            ),
            nn.BatchNorm2d(out_channels, momentum=momentum),
            nn.ReLU(),
        )
        self.conv2 = nn.ModuleList()
        self.conv2.append(ConvBlockRes(out_channels * 2, out_channels, momentum))
        for _ in range(n_blocks - 1):
            self.conv2.append(ConvBlockRes(out_channels, out_channels, momentum))

    def forward(self, x: torch.Tensor, concat_tensor: torch.Tensor) -> torch.Tensor:
        x = torch.cat((self.conv1(x), concat_tensor), dim=1)
        for conv in self.conv2:
            x = conv(x)
        return x


class Decoder(nn.Module):
    def __init__(self, in_channels, n_decoders, stride, n_blocks, momentum=0.01):
        super().__init__()
        self.layers = nn.ModuleList()
        self.n_decoders = n_decoders
        for _ in range(self.n_decoders):
            out_channels = in_channels // 2
            self.layers.append(ResDecoderBlock(in_channels, out_channels, stride, n_blocks, momentum))
            in_channels = out_channels

    def forward(self, x: torch.Tensor, concat_tensors: list) -> torch.Tensor:
        for i in range(self.n_decoders):
            x = self.layers[i](x, concat_tensors[-1 - i])
        return x


class DeepUnet(nn.Module):
    def __init__(
        self,
        kernel_size,
        n_blocks: int,
        en_de_layers: int = 5,
        inter_layers: int = 4,
        in_channels: int = 1,
        en_out_channels: int = 16,
    ):
        super().__init__()
        self.encoder = Encoder(in_channels, 128, en_de_layers, kernel_size, n_blocks, en_out_channels)
        self.intermediate = Intermediate(
            self.encoder.out_channel // 2, self.encoder.out_channel, inter_layers, n_blocks
        )
        self.decoder = Decoder(self.encoder.out_channel, en_de_layers, kernel_size, n_blocks)

    def forward(self, x: torch.Tensor) -> torch.Tensor:
        x, concat_tensors = self.encoder(x)
        x = self.intermediate(x)
        return self.decoder(x, concat_tensors)


class E2E(nn.Module):
    def __init__(self, n_blocks: int, n_gru: int, kernel_size, en_out_channels: int = 16):
        super().__init__()
        self.unet = DeepUnet(kernel_size, n_blocks, 5, 4, 1, en_out_channels)
        self.cnn = nn.Conv2d(en_out_channels, 3, (3, 3), padding=1)
        if n_gru:
            self.fc = nn.Sequential(
                BiGRU(3 * 128, 256, n_gru),
                nn.Linear(512, 360),
                nn.Dropout(0.25),
                nn.Sigmoid(),
            )

    def forward(self, mel: torch.Tensor) -> torch.Tensor:
        mel = mel.transpose(-1, -2).unsqueeze(1)
        x = self.cnn(self.unet(mel)).transpose(1, 2).flatten(-2)
        return self.fc(x)


class MelSpectrogram(nn.Module):
    def __init__(self, is_half: bool, n_mel_channels: int = 128, sample_rate: int = 16000):
        super().__init__()
        self.n_mel_channels = n_mel_channels
        self.sample_rate = sample_rate
        self.is_half = is_half
        self.clamp = 1e-5
        mel_basis = librosa_mel(
            sr=sample_rate, n_fft=1024, n_mels=n_mel_channels, fmin=30, fmax=8000
        )
        self.register_buffer("mel_basis", torch.from_numpy(mel_basis).float())

    def forward(self, audio: torch.Tensor) -> torch.Tensor:
        if audio.dim() == 1:
            audio = audio.unsqueeze(0)
        fft = torch.stft(
            audio,
            n_fft=1024,
            hop_length=160,
            win_length=1024,
            window=torch.hann_window(1024).to(audio.device),
            center=True,
            return_complex=True,
        )
        magnitude = torch.sqrt(fft.real.pow(2) + fft.imag.pow(2))
        mel_output = torch.matmul(self.mel_basis.to(audio.device), magnitude)
        if self.is_half:
            mel_output = mel_output.half()
        return torch.log(torch.clamp(mel_output, min=self.clamp))


class RMVPE:
    hop_length = 160
    sample_rate = 16000

    def __init__(self, model_path: str, is_half: bool, device: str):
        self.is_half = is_half
        self.device = torch.device(device)
        self.mel_extractor = MelSpectrogram(is_half).to(self.device)
        model = E2E(4, 1, (2, 2))
        ckpt = torch.load(model_path, map_location="cpu", weights_only=False)
        if isinstance(ckpt, dict) and "model" in ckpt and not any(k.startswith("unet.") for k in ckpt):
            ckpt = ckpt["model"]
        model.load_state_dict(ckpt)
        model.eval()
        model = model.half() if is_half else model.float()
        self.model = model.to(self.device)
        cents_mapping = 20 * np.arange(360) + 1997.3794084376191
        self.cents_mapping = np.pad(cents_mapping, (4, 4))

    def infer_from_audio(self, audio: np.ndarray, thred: float = 0.03) -> np.ndarray:
        if not torch.is_tensor(audio):
            audio_t = torch.from_numpy(np.asarray(audio, dtype=np.float32))
        else:
            audio_t = audio.float()
        with torch.no_grad():
            mel = self.mel_extractor(audio_t.to(self.device).unsqueeze(0))
            n_frames = mel.shape[-1]
            n_pad = 32 * ((n_frames - 1) // 32 + 1) - n_frames
            if n_pad > 0:
                mel = F.pad(mel, (0, n_pad), mode="constant")
            mel = mel.half() if self.is_half else mel.float()
            hidden = self.model(mel)[:, :n_frames]
        hidden_np = hidden.squeeze(0).float().cpu().numpy()
        return self._decode(hidden_np, thred=thred)

    def _decode(self, hidden: np.ndarray, thred: float) -> np.ndarray:
        cents_pred = self._to_local_average_cents(hidden, thred=thred)
        f0 = 10 * (2 ** (cents_pred / 1200))
        f0[f0 == 10] = 0
        return f0

    def _to_local_average_cents(self, salience: np.ndarray, thred: float) -> np.ndarray:
        center = np.argmax(salience, axis=1)
        salience = np.pad(salience, ((0, 0), (4, 4)))
        center = center + 4
        todo_salience = []
        todo_cents = []
        for idx, c in enumerate(center):
            todo_salience.append(salience[idx, c - 4 : c + 5])
            todo_cents.append(self.cents_mapping[c - 4 : c + 5])
        todo_salience_a = np.asarray(todo_salience)
        todo_cents_a = np.asarray(todo_cents)
        divided = np.sum(todo_salience_a * todo_cents_a, 1) / np.maximum(np.sum(todo_salience_a, 1), 1e-8)
        maxx = np.max(salience, axis=1)
        divided[maxx <= thred] = 0
        return divided
