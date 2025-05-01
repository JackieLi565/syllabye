resource "aws_ses_domain_identity" "this" {
  count  = var.is_dev ? 1 : 0
  domain = var.domain
}

resource "aws_ses_template" "welcome" {
  name    = var.welcome_template_name
  subject = "Welcome to Syllabye, {{name}}!"
  text    = "Hi {{name}}, thanks for joining Syllabye. Your account is now active."
  html    = <<EOT
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <title>Welcome to Syllabye</title>
    <style media="all" type="text/css">
      @media all {
        .btn-primary table td:hover {
          background-color: #ec0867 !important;
        }

        .btn-primary a:hover {
          background-color: #ec0867 !important;
          border-color: #ec0867 !important;
        }
      }
      @media only screen and (max-width: 640px) {
        .main p,
        .main td,
        .main span {
          font-size: 16px !important;
        }

        .wrapper {
          padding: 8px !important;
        }

        .content {
          padding: 0 !important;
        }

        .container {
          padding: 0 !important;
          padding-top: 8px !important;
          width: 100% !important;
        }

        .main {
          border-left-width: 0 !important;
          border-radius: 0 !important;
          border-right-width: 0 !important;
        }

        .btn table {
          max-width: 100% !important;
          width: 100% !important;
        }

        .btn a {
          font-size: 16px !important;
          max-width: 100% !important;
          width: 100% !important;
        }
      }
      @media all {
        .ExternalClass {
          width: 100%;
        }

        .ExternalClass,
        .ExternalClass p,
        .ExternalClass span,
        .ExternalClass font,
        .ExternalClass td,
        .ExternalClass div {
          line-height: 100%;
        }

        .apple-link a {
          color: inherit !important;
          font-family: inherit !important;
          font-size: inherit !important;
          font-weight: inherit !important;
          line-height: inherit !important;
          text-decoration: none !important;
        }

        #MessageViewBody a {
          color: inherit;
          text-decoration: none;
          font-size: inherit;
          font-family: inherit;
          font-weight: inherit;
          line-height: inherit;
        }
      }
    </style>
  </head>
  <body
    style="
      font-family: Helvetica, sans-serif;
      -webkit-font-smoothing: antialiased;
      font-size: 16px;
      line-height: 1.3;
      -ms-text-size-adjust: 100%;
      -webkit-text-size-adjust: 100%;
      background-color: #f4f5f6;
      margin: 0;
      padding: 0;
    "
  >
    <table
      role="presentation"
      border="0"
      cellpadding="0"
      cellspacing="0"
      class="body"
      style="
        border-collapse: separate;
        mso-table-lspace: 0pt;
        mso-table-rspace: 0pt;
        background-color: #f4f5f6;
        width: 100%;
      "
      width="100%"
      bgcolor="#f4f5f6"
    >
      <tr>
        <td
          style="
            font-family: Helvetica, sans-serif;
            font-size: 16px;
            vertical-align: top;
          "
          valign="top"
        >
          &nbsp;
        </td>
        <td
          class="container"
          style="
            font-family: Helvetica, sans-serif;
            font-size: 16px;
            vertical-align: top;
            max-width: 600px;
            padding: 0;
            padding-top: 24px;
            width: 600px;
            margin: 0 auto;
          "
          width="600"
          valign="top"
        >
          <div
            class="content"
            style="
              box-sizing: border-box;
              display: block;
              margin: 0 auto;
              max-width: 600px;
              padding: 0;
            "
          >
            <table
              role="presentation"
              border="0"
              cellpadding="0"
              cellspacing="0"
              class="main"
              style="
                border-collapse: separate;
                mso-table-lspace: 0pt;
                mso-table-rspace: 0pt;
                background: #ffffff;
                border: 1px solid #eaebed;
                border-radius: 16px;
                width: 100%;
              "
              width="100%"
            >
              <tr>
                <td
                  class="wrapper"
                  style="
                    font-family: Helvetica, sans-serif;
                    font-size: 16px;
                    vertical-align: top;
                    box-sizing: border-box;
                    padding: 24px;
                  "
                  valign="top"
                >
                  <p
                    style="
                      font-family: Helvetica, sans-serif;
                      font-size: 16px;
                      font-weight: normal;
                      margin: 0;
                      margin-bottom: 16px;
                    "
                  >
                    Hey {{name}},
                  </p>
                  <p
                    style="
                      font-family: Helvetica, sans-serif;
                      font-size: 16px;
                      font-weight: normal;
                      margin: 0;
                      margin-bottom: 16px;
                    "
                  >
                    Thanks for joining Syllabye! You are receiving this email as
                    a confirmation that your account has been successfully
                    registered.
                  </p>
                  <p
                    style="
                      font-family: Helvetica, sans-serif;
                      font-size: 16px;
                      font-weight: bold;
                      margin: 0;
                      margin-bottom: 16px;
                    "
                  >
                    Questions?
                  </p>
                  <p
                    style="
                      font-family: Helvetica, sans-serif;
                      font-size: 16px;
                      font-weight: normal;
                      margin: 0;
                      margin-bottom: 16px;
                    "
                  >
                    Feel free to reach out to us at
                    <span style="text-decoration: underline"
                      >TODO@torontomu.ca</span
                    >
                    at any time.
                  </p>
                </td>
              </tr>
            </table>

            <!-- START FOOTER -->
            <div
              class="footer"
              style="
                clear: both;
                padding-top: 24px;
                text-align: center;
                width: 100%;
              "
            >
              <table
                role="presentation"
                border="0"
                cellpadding="0"
                cellspacing="0"
                style="
                  border-collapse: separate;
                  mso-table-lspace: 0pt;
                  mso-table-rspace: 0pt;
                  width: 100%;
                "
                width="100%"
              >
                <tr>
                  <td
                    class="content-block"
                    style="
                      font-family: Helvetica, sans-serif;
                      vertical-align: top;
                      color: #9a9ea6;
                      font-size: 16px;
                      text-align: center;
                    "
                    valign="top"
                    align="center"
                  >
                    <span
                      class="apple-link"
                      style="
                        color: #9a9ea6;
                        font-size: 16px;
                        text-align: center;
                      "
                      >Syllabye Co.</span
                    >
                    <br />
                    Don't like these emails?
                    <a
                      href="https://syllabye.cs/unsubscribe"
                      style="
                        color: #9a9ea6;
                        font-size: 16px;
                        text-align: center;
                        text-decoration: none;
                      "
                      >Unsubscribe</a
                    >.
                  </td>
                </tr>
                <tr>
                  <td
                    class="content-block powered-by"
                    style="
                      font-family: Helvetica, sans-serif;
                      vertical-align: top;
                      color: #9a9ea6;
                      font-size: 16px;
                      text-align: center;
                    "
                    valign="top"
                    align="center"
                  >
                    Powered by
                    <a
                      href="https://aws.amazon.com/ses/"
                      style="
                        color: #9a9ea6;
                        font-size: 16px;
                        text-align: center;
                        text-decoration: none;
                      "
                      >Amazon Web Services</a
                    >
                  </td>
                </tr>
              </table>
            </div>
          </div>
        </td>
        <td
          style="
            font-family: Helvetica, sans-serif;
            font-size: 16px;
            vertical-align: top;
          "
          valign="top"
        >
          &nbsp;
        </td>
      </tr>
    </table>
  </body>
</html>
EOT
}

resource "aws_ses_template" "upload_success" {
  name    = var.upload_success_template_name
  subject = "Syllabus Submission Receipt"
  text    = "Hey {{name}} your syllabus submission for {{course}} was not completed successfully."
  html    = <<EOT
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <title>Syllabus Upload Successful</title>
    <style media="all" type="text/css">
      @media all {
        .btn-primary table td:hover {
          background-color: #ec0867 !important;
        }

        .btn-primary a:hover {
          background-color: #ec0867 !important;
          border-color: #ec0867 !important;
        }
      }
      @media only screen and (max-width: 640px) {
        .main p,
        .main td,
        .main span {
          font-size: 16px !important;
        }

        .wrapper {
          padding: 8px !important;
        }

        .content {
          padding: 0 !important;
        }

        .container {
          padding: 0 !important;
          padding-top: 8px !important;
          width: 100% !important;
        }

        .main {
          border-left-width: 0 !important;
          border-radius: 0 !important;
          border-right-width: 0 !important;
        }

        .btn table {
          max-width: 100% !important;
          width: 100% !important;
        }

        .btn a {
          font-size: 16px !important;
          max-width: 100% !important;
          width: 100% !important;
        }
      }
      @media all {
        .ExternalClass {
          width: 100%;
        }

        .ExternalClass,
        .ExternalClass p,
        .ExternalClass span,
        .ExternalClass font,
        .ExternalClass td,
        .ExternalClass div {
          line-height: 100%;
        }

        .apple-link a {
          color: inherit !important;
          font-family: inherit !important;
          font-size: inherit !important;
          font-weight: inherit !important;
          line-height: inherit !important;
          text-decoration: none !important;
        }

        #MessageViewBody a {
          color: inherit;
          text-decoration: none;
          font-size: inherit;
          font-family: inherit;
          font-weight: inherit;
          line-height: inherit;
        }
      }
    </style>
  </head>
  <body
    style="
      font-family: Helvetica, sans-serif;
      -webkit-font-smoothing: antialiased;
      font-size: 16px;
      line-height: 1.3;
      -ms-text-size-adjust: 100%;
      -webkit-text-size-adjust: 100%;
      background-color: #f4f5f6;
      margin: 0;
      padding: 0;
    "
  >
    <table
      role="presentation"
      border="0"
      cellpadding="0"
      cellspacing="0"
      class="body"
      style="
        border-collapse: separate;
        mso-table-lspace: 0pt;
        mso-table-rspace: 0pt;
        background-color: #f4f5f6;
        width: 100%;
      "
      width="100%"
      bgcolor="#f4f5f6"
    >
      <tr>
        <td
          style="
            font-family: Helvetica, sans-serif;
            font-size: 16px;
            vertical-align: top;
          "
          valign="top"
        >
          &nbsp;
        </td>
        <td
          class="container"
          style="
            font-family: Helvetica, sans-serif;
            font-size: 16px;
            vertical-align: top;
            max-width: 600px;
            padding: 0;
            padding-top: 24px;
            width: 600px;
            margin: 0 auto;
          "
          width="600"
          valign="top"
        >
          <div
            class="content"
            style="
              box-sizing: border-box;
              display: block;
              margin: 0 auto;
              max-width: 600px;
              padding: 0;
            "
          >
            <table
              role="presentation"
              border="0"
              cellpadding="0"
              cellspacing="0"
              class="main"
              style="
                border-collapse: separate;
                mso-table-lspace: 0pt;
                mso-table-rspace: 0pt;
                background: #ffffff;
                border: 1px solid #eaebed;
                border-radius: 16px;
                width: 100%;
              "
              width="100%"
            >
              <tr>
                <td
                  class="wrapper"
                  style="
                    font-family: Helvetica, sans-serif;
                    font-size: 16px;
                    vertical-align: top;
                    box-sizing: border-box;
                    padding: 24px;
                  "
                  valign="top"
                >
                  <p
                    style="
                      font-family: Helvetica, sans-serif;
                      font-size: 16px;
                      font-weight: normal;
                      margin: 0;
                      margin-bottom: 16px;
                    "
                  >
                    {{name}},
                  </p>
                  <p
                    style="
                      font-family: Helvetica, sans-serif;
                      font-size: 16px;
                      font-weight: normal;
                      margin: 0;
                      margin-bottom: 16px;
                    "
                  >
                    This email is to confirm that your syllabus submission for
                    {{course}} has been successful.
                  </p>
                  <p
                    style="
                      font-family: Helvetica, sans-serif;
                      font-size: 16px;
                      font-weight: normal;
                      margin: 0;
                      margin-bottom: 16px;
                    "
                  >
                    Thank you for contributing to the Syllabye community — your
                    submission helps other students make informed course
                    decisions.
                  </p>
                  <p
                    style="
                      font-family: Helvetica, sans-serif;
                      font-size: 16px;
                      font-weight: normal;
                      margin: 0;
                      margin-bottom: 16px;
                    "
                  >
                    If you have any questions or need to update your submission,
                    feel free to contact us at
                    <span style="text-decoration: underline; font-weight: bold"
                      >TODO@torontomu.ca</span
                    >
                    at any time.
                  </p>

                  The Syllabye Team
                </td>
              </tr>
            </table>

            <div
              class="footer"
              style="
                clear: both;
                padding-top: 24px;
                text-align: center;
                width: 100%;
              "
            >
              <table
                role="presentation"
                border="0"
                cellpadding="0"
                cellspacing="0"
                style="
                  border-collapse: separate;
                  mso-table-lspace: 0pt;
                  mso-table-rspace: 0pt;
                  width: 100%;
                "
                width="100%"
              >
                <tr>
                  <td
                    class="content-block"
                    style="
                      font-family: Helvetica, sans-serif;
                      vertical-align: top;
                      color: #9a9ea6;
                      font-size: 16px;
                      text-align: center;
                    "
                    valign="top"
                    align="center"
                  >
                    <span
                      class="apple-link"
                      style="
                        color: #9a9ea6;
                        font-size: 16px;
                        text-align: center;
                      "
                      >Syllabye Co.</span
                    >
                    <br />
                  </td>
                </tr>
                <tr>
                  <td
                    class="content-block powered-by"
                    style="
                      font-family: Helvetica, sans-serif;
                      vertical-align: top;
                      color: #9a9ea6;
                      font-size: 16px;
                      text-align: center;
                    "
                    valign="top"
                    align="center"
                  >
                    Powered by
                    <a
                      href="https://aws.amazon.com/ses/"
                      style="
                        color: #9a9ea6;
                        font-size: 16px;
                        text-align: center;
                        text-decoration: none;
                      "
                      >Amazon Web Services</a
                    >
                  </td>
                </tr>
              </table>
            </div>
          </div>
        </td>
        <td
          style="
            font-family: Helvetica, sans-serif;
            font-size: 16px;
            vertical-align: top;
          "
          valign="top"
        >
          &nbsp;
        </td>
      </tr>
    </table>
  </body>
</html>
EOT
}

resource "aws_ses_template" "upload_error" {
  name    = var.upload_error_template_name
  subject = "Syllabus Submission Error"
  html    = <<EOT
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <title>Syllabus Upload Unsuccessful</title>
    <style media="all" type="text/css">
      @media all {
        .btn-primary table td:hover {
          background-color: #ec0867 !important;
        }

        .btn-primary a:hover {
          background-color: #ec0867 !important;
          border-color: #ec0867 !important;
        }
      }
      @media only screen and (max-width: 640px) {
        .main p,
        .main td,
        .main span {
          font-size: 16px !important;
        }

        .wrapper {
          padding: 8px !important;
        }

        .content {
          padding: 0 !important;
        }

        .container {
          padding: 0 !important;
          padding-top: 8px !important;
          width: 100% !important;
        }

        .main {
          border-left-width: 0 !important;
          border-radius: 0 !important;
          border-right-width: 0 !important;
        }

        .btn table {
          max-width: 100% !important;
          width: 100% !important;
        }

        .btn a {
          font-size: 16px !important;
          max-width: 100% !important;
          width: 100% !important;
        }
      }
      @media all {
        .ExternalClass {
          width: 100%;
        }

        .ExternalClass,
        .ExternalClass p,
        .ExternalClass span,
        .ExternalClass font,
        .ExternalClass td,
        .ExternalClass div {
          line-height: 100%;
        }

        .apple-link a {
          color: inherit !important;
          font-family: inherit !important;
          font-size: inherit !important;
          font-weight: inherit !important;
          line-height: inherit !important;
          text-decoration: none !important;
        }

        #MessageViewBody a {
          color: inherit;
          text-decoration: none;
          font-size: inherit;
          font-family: inherit;
          font-weight: inherit;
          line-height: inherit;
        }
      }
    </style>
  </head>
  <body
    style="
      font-family: Helvetica, sans-serif;
      -webkit-font-smoothing: antialiased;
      font-size: 16px;
      line-height: 1.3;
      -ms-text-size-adjust: 100%;
      -webkit-text-size-adjust: 100%;
      background-color: #f4f5f6;
      margin: 0;
      padding: 0;
    "
  >
    <table
      role="presentation"
      border="0"
      cellpadding="0"
      cellspacing="0"
      class="body"
      style="
        border-collapse: separate;
        mso-table-lspace: 0pt;
        mso-table-rspace: 0pt;
        background-color: #f4f5f6;
        width: 100%;
      "
      width="100%"
      bgcolor="#f4f5f6"
    >
      <tr>
        <td
          style="
            font-family: Helvetica, sans-serif;
            font-size: 16px;
            vertical-align: top;
          "
          valign="top"
        >
          &nbsp;
        </td>
        <td
          class="container"
          style="
            font-family: Helvetica, sans-serif;
            font-size: 16px;
            vertical-align: top;
            max-width: 600px;
            padding: 0;
            padding-top: 24px;
            width: 600px;
            margin: 0 auto;
          "
          width="600"
          valign="top"
        >
          <div
            class="content"
            style="
              box-sizing: border-box;
              display: block;
              margin: 0 auto;
              max-width: 600px;
              padding: 0;
            "
          >
            <table
              role="presentation"
              border="0"
              cellpadding="0"
              cellspacing="0"
              class="main"
              style="
                border-collapse: separate;
                mso-table-lspace: 0pt;
                mso-table-rspace: 0pt;
                background: #ffffff;
                border: 1px solid #eaebed;
                border-radius: 16px;
                width: 100%;
              "
              width="100%"
            >
              <tr>
                <td
                  class="wrapper"
                  style="
                    font-family: Helvetica, sans-serif;
                    font-size: 16px;
                    vertical-align: top;
                    box-sizing: border-box;
                    padding: 24px;
                  "
                  valign="top"
                >
                  <p
                    style="
                      font-family: Helvetica, sans-serif;
                      font-size: 16px;
                      font-weight: normal;
                      margin: 0;
                      margin-bottom: 16px;
                    "
                  >
                    {{name}},
                  </p>
                  <p
                    style="
                      font-family: Helvetica, sans-serif;
                      font-size: 16px;
                      font-weight: normal;
                      margin: 0;
                      margin-bottom: 16px;
                    "
                  >
                    Unfortunately, your syllabus submission for {{course}} was
                    not completed successfully.
                  </p>

                  <p
                    style="
                      font-family: Helvetica, sans-serif;
                      font-size: 16px;
                      font-weight: normal;
                      margin: 0;
                      margin-bottom: 16px;
                    "
                  >
                    {{reason}}
                  </p>
                  <p
                    style="
                      font-family: Helvetica, sans-serif;
                      font-size: 16px;
                      font-weight: normal;
                      margin: 0;
                      margin-bottom: 16px;
                    "
                  >
                    If this does&apos;t seem right please try again. If the
                    problem continues, feel free to reach out to us at
                    <span style="text-decoration: underline; font-weight: bold"
                      >TODO@torontomu.ca</span
                    >
                  </p>
                  <p
                    style="
                      font-family: Helvetica, sans-serif;
                      font-size: 16px;
                      font-weight: normal;
                      margin: 0;
                      margin-bottom: 16px;
                    "
                  >
                    Thank you for your contribution to the Syllabye community.
                  </p>

                  The Syllabye Team
                </td>
              </tr>
            </table>

            <div
              class="footer"
              style="
                clear: both;
                padding-top: 24px;
                text-align: center;
                width: 100%;
              "
            >
              <table
                role="presentation"
                border="0"
                cellpadding="0"
                cellspacing="0"
                style="
                  border-collapse: separate;
                  mso-table-lspace: 0pt;
                  mso-table-rspace: 0pt;
                  width: 100%;
                "
                width="100%"
              >
                <tr>
                  <td
                    class="content-block"
                    style="
                      font-family: Helvetica, sans-serif;
                      vertical-align: top;
                      color: #9a9ea6;
                      font-size: 16px;
                      text-align: center;
                    "
                    valign="top"
                    align="center"
                  >
                    <span
                      class="apple-link"
                      style="
                        color: #9a9ea6;
                        font-size: 16px;
                        text-align: center;
                      "
                      >Syllabye Co.</span
                    >
                    <br />
                  </td>
                </tr>
                <tr>
                  <td
                    class="content-block powered-by"
                    style="
                      font-family: Helvetica, sans-serif;
                      vertical-align: top;
                      color: #9a9ea6;
                      font-size: 16px;
                      text-align: center;
                    "
                    valign="top"
                    align="center"
                  >
                    Powered by
                    <a
                      href="https://aws.amazon.com/ses/"
                      style="
                        color: #9a9ea6;
                        font-size: 16px;
                        text-align: center;
                        text-decoration: none;
                      "
                      >Amazon Web Services</a
                    >
                  </td>
                </tr>
              </table>
            </div>
          </div>
        </td>
        <td
          style="
            font-family: Helvetica, sans-serif;
            font-size: 16px;
            vertical-align: top;
          "
          valign="top"
        >
          &nbsp;
        </td>
      </tr>
    </table>
  </body>
</html>
EOT
}
