from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
import os
import shutil
from utils.image_loader import process_drive_pdf
from stages.ocr import extract_text
from stages.phash import process_phash_for_image
from stages.pdf_name_forensics import run_pdf_name_forensics
from stages.cnn_infer_anomaly import run_cnn_anomaly
from stages.aggregator import aggregate_verdict

app = FastAPI()

class VerifyRequest(BaseModel):
    drive_link: str

@app.post("/verify")
def verify_certificate(req: VerifyRequest):
    print(f"🚀 Verifying: {req.drive_link}")

    # 0. Forensics

    try:
        pdf_forensics = run_pdf_name_forensics(req.drive_link)
    except Exception as e:
        print(f"⚠️ Forensics/Download failed: {e}")
        return {
            "trust_score": 0.0,
            "final_verdict": "SUSPICIOUS",
            "components": {
                "forensics": {"error": str(e)},
                "phash": {},
                "cnn": {}
            }
        }

    # 1. Download & Convert
    try:
        image_paths = process_drive_pdf(req.drive_link)
    except Exception as e:
        raise HTTPException(status_code=400, detail=f"Failed to process PDF: {str(e)}")

    if not image_paths:
        raise HTTPException(status_code=400, detail="No images generated from drive link")

    image_path = os.path.abspath(image_paths[0])
    
    # 2. OCR
    ocr_out = extract_text(image_path)
    if not ocr_out["is_text_found"]:
        # If no text found, we can fail or return low score
        # For now, let's proceed but note it
        pass

    # 3. pHash
    phash_result = process_phash_for_image(
        image_path=image_path,
        ocr_text=ocr_out["raw_text"]
    )

    # 4. CNN
    cnn_result = run_cnn_anomaly(image_path)

    # 5. Aggregate
    final_result = aggregate_verdict(
        pdf_forensics,
        phash_result,
        cnn_result
    )

    # Cleanup image if needed (optional)
    # os.remove(image_path)
    
    return final_result

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
