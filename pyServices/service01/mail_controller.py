import smtplib, ssl, email.encoders
from email.mime.text import MIMEText
from email.mime.base import MIMEBase
from email.mime.image import MIMEImage

from email.mime.multipart import MIMEMultipart


class MailEngine:

    def __init__(self, subject, text):
        self.email_user = "unilever.belen.dev@gmail.com"
        self.email_pass = "belen2718"
        self.fileName = ""
        self.subject =  '"' +  str(subject) + '"'
        self.text = text


    def start(self):

        # to = ['bryan.romero@unilever.com', 'r0ch1bryan@gmail.com', self.email_user]
        to = [self.email_user]

        sender_email = self.email_user
        receiver_email = ','.join(to)
        password = self.email_pass
        text = self.text

        message = MIMEMultipart("alternative")
        message["Subject"] = self.subject
        message["From"] = sender_email
        message["To"] = receiver_email

        Images = ['unilever.png', 'plot.png', 'am_plot.png', 'oee_plot.png']


        html = text


        # Turn these into plain/html MIMEText objects
        part1 = MIMEText(html, "html")

        # Add HTML/plain-text parts to MIMEMultipart message
        # The email client will try to render the last part first
        message.attach(part1)

        i = 0
        for picture in Images:

            # This example assumes the image is in the current directory
            fp = open(picture, 'rb')
            msgImage = MIMEImage(fp.read())
            fp.close()

            # Define the image's ID as referenced above
            # print('<image'+str(i)+'>')
            msgImage.add_header('Content-ID', '<image'+str(i)+'>')
            message.attach(msgImage)
            i += 1


        # Create secure connection with server and send email
        context = ssl.create_default_context()
        with smtplib.SMTP_SSL("smtp.gmail.com", 465, context=context) as server:
            server.login(sender_email, password)
            server.sendmail(
                sender_email, receiver_email, message.as_string()
            )
            server.quit()


if __name__ == "__main__":

    print("EMAIL ENGINE")


    engine = MailEngine("Actividad en Planta Belen", "<b>Test</b>")

    engine.start()
