import boto3
import io
import fitz
import os
from PIL import Image

s3 = boto3.client("s3")


def generate_pdf_thumbnail(pdf_stream, max_size=(200, 260), dpi=72):
    """Converts first page of PDF to a padded thumbnail JPEG in memory."""
    doc = fitz.open(stream=pdf_stream, filetype="pdf")
    page = doc.load_page(0)
    pix = page.get_pixmap(dpi=dpi)

    image = Image.frombytes("RGB", [pix.width, pix.height], pix.samples)
    image.thumbnail(max_size)

    # White background image
    canvas = Image.new("RGB", max_size, (255, 255, 255))
    x_offset = (max_size[0] - image.width) // 2
    y_offset = (max_size[1] - image.height) // 2
    canvas.paste(image, (x_offset, y_offset))

    output = io.BytesIO()
    canvas.save(output, format="JPEG")
    output.seek(0)
    return output


def handler(event, _):
    dest_bucket = os.getenv("AWS_S3_THUMBNAIL_BUCKET")

    try:
        record = event["Records"][0]
        src_bucket = record["s3"]["bucket"]["name"]
        key = record["s3"]["object"]["key"]

        # Retrieve syllabus from bucket
        pdf_stream = io.BytesIO()
        s3.download_fileobj(src_bucket, key, pdf_stream)
        pdf_stream.seek(0)

        thumb_buffer = generate_pdf_thumbnail(pdf_stream)

        dest_key = key + ".jpg"
        s3.upload_fileobj(
            thumb_buffer, dest_bucket, dest_key, ExtraArgs={"ContentType": "image/jpeg"}
        )
        print(f"syllabus thumbnail created: {key}")
    except Exception as e:
        print(f"error processing record: {e}")
