import open_clip
import torch

# Specify the model and dataset
model_name = 'ViT-B-32'
pretrained = 'laion2b_s34b_b79k'

# This will download the model and cache it locally
model, _, preprocess = open_clip.create_model_and_transforms(model_name, pretrained=pretrained)

# Save locally for later use
torch.save(model.state_dict(), f'{model_name}_{pretrained}.pt')