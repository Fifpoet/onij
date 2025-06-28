// src/api/types/ai.ts
export interface AiChatReq {
    prompt: string;
    model: string;
    image_url: string;
}

export interface AiChatResp {
    code: number;
    message: string;
    answer: string;
}