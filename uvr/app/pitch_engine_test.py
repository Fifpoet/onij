import unittest

import numpy as np

from app.pitch_engine import _correct_period_errors, _notes_from_f0, _sticky_merge


class PitchEngineTest(unittest.TestCase):
    def test_sticky_keeps_fifth_jump(self):
        notes = [
            {"t0": 0.0, "t1": 0.2, "midi": 45.0},
            {"t0": 0.21, "t1": 0.40, "midi": 52.0},
        ]
        out = _sticky_merge([dict(n) for n in notes])
        self.assertEqual([n["midi"] for n in out], [45.0, 52.0])

    def test_sticky_smoothes_semitone_blip(self):
        notes = [
            {"t0": 0.0, "t1": 0.2, "midi": 52.0},
            {"t0": 0.21, "t1": 0.28, "midi": 53.0},
            {"t0": 0.29, "t1": 0.5, "midi": 52.0},
        ]
        out = _sticky_merge([dict(n) for n in notes])
        self.assertTrue(all(n["midi"] == 52.0 for n in out))

    def test_correct_period_lifts_fifth_below(self):
        # 主体 E3(52)，混入低五度 A2(45) ≈ 110Hz
        e3 = 164.81
        a2 = 110.0
        f0 = np.array([e3] * 80 + [a2] * 40 + [e3] * 80, dtype=np.float64)
        out = _correct_period_errors(f0)
        midis = 69.0 + 12.0 * np.log2(out / 440.0)
        self.assertGreater(float(np.median(midis[80:120])), 50.0)

    def test_notes_from_stable_e3(self):
        hz = 164.81
        f0 = np.zeros(200, dtype=np.float64)
        f0[20:80] = hz
        notes = _notes_from_f0(f0, 10.0)
        self.assertTrue(notes)
        self.assertEqual(notes[0]["midi"], 52.0)


if __name__ == "__main__":
    unittest.main()
