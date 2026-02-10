import os
import requests
import fitz  # PyMuPDF
from urllib.parse import urlparse, parse_qs

# Poppler path not needed with PyMuPDF

def extract_file_id(drive_link):
    """Extract FILE_ID from any Google Drive link."""
    if "id=" in drive_link:
        return parse_qs(urlparse(drive_link).query)["id"][0]
    elif "/d/" in drive_link:
        return drive_link.split("/d/")[1].split("/")[0]
    else:
        # Fallback for simple ID extraction
        parts = drive_link.split('/')
        for part in parts:
             if len(part) > 25: # simplistic check for ID length
                 return part
        raise ValueError("Invalid Google Drive link")


def download_pdf_temp(drive_link, file_id):
    """Download PDF temporarily (will be deleted after PNG conversion)."""
    # Use direct download link format
    url = f"https://drive.google.com/uc?export=download&id={file_id}"
    temp_pdf = f"temp_{file_id}.pdf"

    print(f"⬇ Downloading PDF ({file_id})...")

    try:
        response = requests.get(url, stream=True)
        response.raise_for_status()

        with open(temp_pdf, "wb") as f:
            for chunk in response.iter_content(chunk_size=8192):
                if chunk:
                    f.write(chunk)
        return temp_pdf
    except Exception as e:
        print(f"Error downloading PDF: {e}")
        # Return none or raise error? Let's raise to be caught by caller
        raise e


def pdf_to_images(pdf_path, file_id, output_folder="output_images"):
    """Convert PDF → PNG using PyMuPDF and DELETE the PDF afterwards."""
    os.makedirs(output_folder, exist_ok=True)

    print("🖼 Converting PDF to PNG using PyMuPDF...")

    output_paths = []
    try:
        doc = fitz.open(pdf_path)
        for i, page in enumerate(doc, start=1):
            # Render page to image (pixmap)
            # matrix=fitz.Matrix(2, 2) roughly doubles resolution, similar to higher DPI
            pix = page.get_pixmap(matrix=fitz.Matrix(2, 2)) 
            out_path = os.path.join(output_folder, f"{file_id}_page_{i}.png")
            pix.save(out_path)
            output_paths.append(out_path)
        
        doc.close()
        print(f"✅ Saved {len(output_paths)} PNG files in {output_folder}")
    except Exception as e:
        print(f"Error converting PDF: {e}")
        raise e

    if os.path.exists(pdf_path):
        try:
            os.remove(pdf_path)
            print(f"🗑 Deleted temporary PDF: {pdf_path}")
        except Exception as e:
             print(f"Warning: Could not delete temp PDF: {e}")

    return output_paths


def process_drive_pdf(drive_link):
    """Drive link → temp PDF → PNG paths"""
    file_id = extract_file_id(drive_link)
    temp_pdf_path = download_pdf_temp(drive_link, file_id)
    png_paths = pdf_to_images(temp_pdf_path, file_id)
    return png_paths


if __name__ == "__main__":
    link = input("Enter Google Drive PDF link: ").strip()
    paths = process_drive_pdf(link)

    print("\n📁 Generated PNG files:")
    for p in paths:
        print(" →", p)
