from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field
import open_clip
import torch
import psycopg2
import psycopg2.extras
import os





device = "cuda" if torch.cuda.is_available() else "cpu"

# loading the model from Memeory   /// pretrained is the destination also
model_name = 'ViT-B-32'
pretrained = 'laion2b_s34b_b79k'
model, _, preprocess = open_clip.create_model_and_transforms(model_name, pretrained=pretrained)
tokenizer = open_clip.get_tokenizer('ViT-B-32')
model = model.to(device)


class InputFeatures(BaseModel):
    title: str = Field(default="")
    description: str  = Field(default="")
    tags: list[str] = Field(default_factory=list)
    user_name: str = Field(default="")
    video_id: int = Field(default=0)


app = FastAPI()

@app.post("/vectorize-video/")
async def VectorEmbedingGenrationVideo(requestData: InputFeatures):
    if requestData.video_id == 0:
        raise HTTPException(status_code=400, detail="Invalid video_id")
        return
    print("Received data:", requestData)
    text_data_join = f"Title: {requestData.title}. Description: {requestData.description}. Tags: {', '.join(requestData.tags)}. Uploaded by {requestData.user_name}."
    text_tokens = tokenizer([text_data_join]).to(device)
    
    with torch.no_grad():
        t_vec = model.encode_text(text_tokens)
        print(t_vec.shape)
    print(text_tokens.shape)

    t_vec /= t_vec.norm(dim=-1, keepdim=True)
    result = t_vec.cpu().numpy().tolist()
    conn =None
    try:
        conn = psycopg2.connect("dbname='MetaDataStorage' user='postgres' host='localhost' password='Narayan!123' port='5454'")
        with conn.cursor(cursor_factory=psycopg2.extras.DictCursor) as curs:
            update_query = """
                UPDATE video_data 
                SET embeddings = %s 
                WHERE video_id = %s;
                """
            curs.execute(update_query, (result[0], requestData.video_id))
            conn.commit()
    except Exception as e:
        print(f"Database operation failed: {e}")
        raise HTTPException(status_code=500, detail="Database operation error")
    finally:
        if conn:
            conn.close()


    return {"vector_embeding": result, "status": "narayan narayan narayan narayan"}


class InputFeaturesUser(BaseModel):
    user_handle: str
    user_profile_name: str
    user_description: str
    from_location: str
    user_date_of_birth: str
    gender: str
    tags: list[str]
    user_id: int   = Field(default=-1)



@app.post("/vectorize-user/")
async def VectorEmbedingGenrationUser(requestData: InputFeaturesUser):
    if requestData.user_id == -1:
        raise HTTPException(status_code=400, detail="Invalid user_id")
        return

    text_data_join = f"Handle: {requestData.user_handle}. Name: {requestData.user_profile_name}. Description: {requestData.user_description}. Location: {requestData.from_location}. DOB: {requestData.user_date_of_birth}. Gender: {requestData.gender}. Tags: {', '.join(requestData.tags)}."
    text_tokens = tokenizer([text_data_join]).to(device)
    
    with torch.no_grad():
        t_vec = model.encode_text(text_tokens)
    
    t_vec /= t_vec.norm(dim=-1, keepdim=True)
    result = t_vec.cpu().numpy().tolist()
    
    conn = None
    try:
        conn = psycopg2.connect("dbname='MetaDataStorage' user='postgres' host='localhost' password='Narayan!123' port='5454'")
        with conn.cursor(cursor_factory=psycopg2.extras.DictCursor) as curs:
            update_query = "UPDATE user_data_table SET embeddings = %s WHERE user_id = %s;"
            curs.execute(update_query, (result[0], requestData.user_id))
            conn.commit()
    except Exception as e:
        raise HTTPException(status_code=500, detail="Database operation error")
    finally:
        if conn:
            conn.close()
    
    return {"vector_embeding": result, "status": "success"}


class InputFeaturesEco(BaseModel):
    eco_text:     		str
    tags:				list[str]
    img_count:		    int
    uploader_name: 		    str
    eco_id:              int = Field(default=0)

@app.post("/vectorize-eco/")
async def VectorEmbedingGenrationEco(requestData: InputFeaturesEco):
    if requestData.EcoId == 0:
        raise HTTPException(status_code=400, detail="Invalid EcoId")
        return
    print("Received data:", requestData)
    text_data_join = f"EcoText: {requestData.eco_text}.  Tags: {', '.join(requestData.tags)}. Uploaded by {requestData.uploader_name}."
    text_tokens = tokenizer([text_data_join]).to(device)
    
    with torch.no_grad():
        t_vec = model.encode_text(text_tokens)
        print(t_vec.shape)
    print(text_tokens.shape)

    t_vec /= t_vec.norm(dim=-1, keepdim=True)
    result = t_vec.cpu().numpy().tolist()
    conn =None
    try:
        conn = psycopg2.connect("dbname='MetaDataStorage' user='postgres' host='localhost' password='Narayan!123' port='5454'")
        with conn.cursor(cursor_factory=psycopg2.extras.DictCursor) as curs:
            update_query = """
                UPDATE eco_data 
                SET embeddings = %s 
                WHERE eco_id = %s;
                """
            curs.execute(update_query, (result[0], requestData.eco_id))
            conn.commit()
    except Exception as e:
        print(f"Database operation failed: {e}")
        raise HTTPException(status_code=500, detail="Database operation error")
    finally:
        if conn:
            conn.close()


    return {"vector_embeding": result, "status": "narayan narayan narayan narayan"}


