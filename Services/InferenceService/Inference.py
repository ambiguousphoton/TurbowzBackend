from fastapi import FastAPI
from pydantic import BaseModel
import open_clip
import torch
device = "cuda" if torch.cuda.is_available() else "cpu"

# loading the model from Memeory   /// pretrained is the destination also
model_name = 'ViT-B-32'
pretrained = 'laion2b_s34b_b79k'
model, _, preprocess = open_clip.create_model_and_transforms(model_name, pretrained=pretrained)
tokenizer = open_clip.get_tokenizer('ViT-B-32')
model = model.to(device)


class InputFeatures(BaseModel):
    title: str
    description: str
    tags: list[str]
    user_name: str


app = FastAPI()

@app.post("/vectorize/")
async def VectorEmbedingGenration(requestData: InputFeatures):
    text_data_join = f"Title: {requestData.title}. Description: {requestData.description}. Tags: {', '.join(requestData.tags)}. Uploaded by {requestData.user_name}."
    text_tokens = tokenizer([text_data_join]).to(device)
    
    with torch.no_grad():
        t_vec = model.encode_text(text_tokens)
        print(t_vec.shape)
    print(text_tokens.shape)

    t_vec /= t_vec.norm(dim=-1, keepdim=True)
    result = t_vec.cpu().numpy().tolist()
    return {"vector_embeding": result, "status": "narayan narayan narayan narayan"}




#  Def  using  transcripts  generated  from   Audio.