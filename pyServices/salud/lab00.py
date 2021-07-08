def write_pdf(self):
        # Write PDF.
        from fpdf import FPDF
        from PIL import Image
        dpi = 120 # note for calcs below that "pt" units are 1/72th of an inch
        pdf = FPDF(unit="pt")
        for image in self.screenshot_image_filenames:
            # Size the next PDF page to the size of this image.
            with open(image, "rb") as f:
                im = Image.open(f)
                page_size = im.size[0]/dpi*72, im.size[1]/dpi*72
                
            pdf.add_page(format=page_size)

            pdf.image(image, 0, 0, page_size[0], page_size[1])

        pdf.output(self.write_pdf_filename, "F")
