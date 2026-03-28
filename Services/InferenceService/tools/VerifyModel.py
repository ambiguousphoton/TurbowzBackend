print("narayan narayan narayan narayan")
import torch
import open_clip

# Define model
model_name = 'ViT-B-32'
pretrained = 'laion2b_s34b_b79k'

# Try loading from cache or local
print("🔍 Loading model...")
model, _, preprocess = open_clip.create_model_and_transforms(model_name, pretrained=pretrained)

# Move to CPU (or GPU if available)
device = "cuda" if torch.cuda.is_available() else "cpu"
model = model.to(device)

# Quick sanity check
print("✅ Model loaded successfully.")
print(f"Model device: {next(model.parameters()).device}")
print(f"Model architecture: {model_name}")
