package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/mailgun/mailgun-go"

	// "strconv"
	"log"
	"net/http"

	// "io/ioutil"
	"os"
	// "errors"
	"strings"
	"time"
)

/*
*************************************
MODEL

	**************************************
*/
type Contact struct {
	Uuid        string `json:"uuid"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone 		string `json:"phone"`
	Project     string `json:"project"`
	Source      string `json:"source"`
	UtmSources  string `json:"utm_sources"`
	Terms       bool   `json:"terms"`
	CreatedAt   string `json:"created_at"`
}

var ssmParams map[string]string

/*
*************************************
DB START

	**************************************
*/
func SaveToDynamoDb(info *Contact) error {
	sess, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		return err
	}

	svc := dynamodb.NewFromConfig(sess)

	uid, _ := uuid.NewUUID()
	loc, _ := time.LoadLocation("Asia/Kuala_Lumpur")
	t := time.Now().In(loc)
	info.Uuid = uid.String()
	info.CreatedAt = t.Format("2006-01-02 15:04:05")

	if info.Source == "GDN" {
		info.Source = "GDN - Google Ads"
	} else if info.Source == "GDNV1" {
		info.Source = "GDN - Google Ads"
	} else if info.Source == "GDNV2" {
		info.Source = "GDNV2"
	} else if info.Source == "FBAds" {
		info.Source = "FBAds"
	} else if info.Source == "FBForm" {
		info.Source = "Facebook Lead Form"
	} else if info.Source == "FBWeb" {
		info.Source = "FBWeb"
	} else if info.Source == "SearchOP1" {
		info.Source = "SearchOP1"
	} else if info.Source == "SearchOP2" {
		info.Source = "SearchOP2"
	} else if info.Source == "LinkedIn" {
		info.Source = "LinkedIn"
	} else if info.Source == "WABlast" {
		info.Source = "WABlast"
	} else if info.Source == "innity" {
		info.Source = "innity"
	} else if info.Source == "GLWebsite" {
		info.Source = "GLWebsite"
	} else if info.Source == "GLOPD" {
		info.Source = "GLOPD"
	} else if info.Source == "GLEDM" {
		info.Source = "GLEDM"
	} else if info.Source == "FBPost" {
		info.Source = "FBPost"
	} else {
		info.Source = "Gamuda Cove Product Landing Page"
	}

	var input *dynamodb.PutItemInput

	itemUnmarshal := map[string]string{
		"uuid":        info.Uuid,
		"name":        info.Name,
		"email":       info.Email,
		"phone":       info.Phone,
		"project":     info.Project,
		"source":      info.Source,
		"utm_sources": info.UtmSources,
		"created_at":  info.CreatedAt,
	}

	item, err := attributevalue.MarshalMap(itemUnmarshal)
	if err != nil {
		panic(err)
	}
	input = &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DDBTABLE_PRODUCTS")),
		Item:      item,
	}

	_, err2 := svc.PutItem(context.TODO(), input)

	return err2
}

/**************************************
DB END
 ***************************************/

/*
*************************************
MAIL START

	**************************************
*/
func SendSimpleMail(mailto string, subject string, html string, text string, sender string, mailfrom string) (string, error) {
	mg := mailgun.NewMailgun(os.Getenv("MG_DOMAIN"), ssmParams["MAILGUN_API_KEY"])
	m := mg.NewMessage(
		sender+"<"+mailfrom+">",
		subject,
		text,
		mailto,
	)
	m.SetHtml(html)
	m.SetTracking(true)
	m.SetTrackingClicks(true)
	m.SetTrackingOpens(true)
	//m.SetReplyTo(os.Getenv("MG_REPLY_TO"))

	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	//defer cancel()

	_, id, err := mg.Send(m)

	return id, err
}

/**************************************
MAIL END
 ***************************************/

/*
*************************************
SENTINO START

	**************************************
*/
func SaveToSentino(info *Contact) error {
	url := os.Getenv("SENTINO_ENDPOINT")

	if info.Source == "GDN - Google Ads" {
		os.Setenv("SENTINO_SOURCE_ID", "1217i8Ev1hT20190517154848")
	} else if info.Source == "GDNV1" {
		os.Setenv("SENTINO_SOURCE_ID", "1217i8Ev1hT20190517154848")
	} else if info.Source == "GDNV2" {
		os.Setenv("SENTINO_SOURCE_ID", "1519Jxnwc5U20191213164309")
	} else if info.Source == "FBAds" {
		os.Setenv("SENTINO_SOURCE_ID", "170POu0ZneH20200212184627")
	} else if info.Source == "Facebook Lead Form" {
		os.Setenv("SENTINO_SOURCE_ID", "122CYsyXyF820190517155007")
	} else if info.Source == "FBWeb" {
		os.Setenv("SENTINO_SOURCE_ID", "171iIjaUkwK20200212190451")
	} else if info.Source == "SearchOP1" {
		os.Setenv("SENTINO_SOURCE_ID", "163iIZiPsGz20200124204449")
	} else if info.Source == "SearchOP2" {
		os.Setenv("SENTINO_SOURCE_ID", "164NbcSkzi120200124204517")
	} else if info.Source == "LinkedIn" {
		os.Setenv("SENTINO_SOURCE_ID", "173hqClzTbG20200212190726")
	} else if info.Source == "WABlast" {
		os.Setenv("SENTINO_SOURCE_ID", "178LegaLjEd20200219110500")
	} else if info.Source == "innity" {
		os.Setenv("SENTINO_SOURCE_ID", "131jHnoX7Qa20190816170029")
	} else if info.Source == "GLWebsite" {
		os.Setenv("SENTINO_SOURCE_ID", "189hn13LYM720200322221142")
	} else if info.Source == "GLOPD" {
		os.Setenv("SENTINO_SOURCE_ID", "188FpqwK2ww20200322033424")
	} else if info.Source == "GLEDM" {
		os.Setenv("SENTINO_SOURCE_ID", "20671SPBABA20200427142516")
	} else if info.Source == "FBPost" {
		os.Setenv("SENTINO_SOURCE_ID", "160HtxnTgZE20200123230808")
	} else {
		os.Setenv("SENTINO_SOURCE_ID", "450mG0qaN8w20210513180301")
	}

	payload := []byte(strings.TrimSpace(`
		<soapenv:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:urn="urn:Service2Controllerwsdl">
   			<soapenv:Header/>
   			<soapenv:Body>
				<urn:setRegistration soapenv:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
					<accessToken xsi:type="xsd:string">` + ssmParams["SENTINO_ACCESS_TOKEN"] + `</accessToken>
         			<fname xsi:type="xsd:string">` + info.Name + `</fname>
					<lname xsi:type="xsd:string"></lname>
         			<email xsi:type="xsd:string">` + info.Email + `</email>
         			<mobile xsi:type="xsd:string">` + info.Phone + `</mobile>
         			<office_number xsi:type="xsd:string"></office_number>
         			<ic_number xsi:type="xsd:string"></ic_number>
         			<address xsi:type="xsd:string"></address>
         			<city xsi:type="xsd:string"></city>
         			<postcode xsi:type="xsd:string"></postcode>
         			<state xsi:type="xsd:string"></state>
         			<bumi_status xsi:type="xsd:string"></bumi_status>
         			<industry xsi:type="xsd:string"></industry>
         			<race xsi:type="xsd:string"></race>
         			<age xsi:type="xsd:string"></age>
         			<age_group xsi:type="xsd:string"></age_group>
         			<income_bracket xsi:type="xsd:string"></income_bracket>
         			<income xsi:type="xsd:string"></income>
         			<gender xsi:type="xsd:string"></gender>
         			<date_of_birth xsi:type="xsd:string"></date_of_birth>
         			<occupation xsi:type="xsd:string"></occupation>
         			<nationality xsi:type="xsd:string"></nationality>
         			<language xsi:type="xsd:string"></language>
         			<family_size xsi:type="xsd:string"></family_size>
         			<buying_reason xsi:type="xsd:string"></buying_reason>
         			<propertyType xsi:type="xsd:string"></propertyType>
         			<preferred_location xsi:type="xsd:string"></preferred_location>
         			<preferred_state xsi:type="xsd:string"></preferred_state>
         			<preferred_price_range xsi:type="xsd:string"></preferred_price_range>
         			<preferred_price xsi:type="xsd:string"></preferred_price>
         			<furnishing xsi:type="xsd:string"></furnishing>
         			<min_rooms xsi:type="xsd:string"></min_rooms>
         			<property_siz xsi:type="xsd:string"></property_siz>
         			<facilities xsi:type="xsd:string"></facilities>
         			<accessibility xsi:type="xsd:string"></accessibility>
         			<amenities xsi:type="xsd:string"></amenities>
         			<completion_date xsi:type="xsd:string"></completion_date>
         			<development_stage xsi:type="xsd:string"></development_stage>
         			<direction xsi:type="xsd:string"></direction>
         			<project xsi:type="xsd:string">` + os.Getenv("SENTINO_PROJECT_ID") + `</project>
         			<source xsi:type="xsd:string">` + os.Getenv("SENTINO_SOURCE_ID") + `</source>
					<remarks xsi:type="xsd:string">Project: ` + info.Project + `</remarks>
					<utm_sources xsi:type="xsd:string">` + info.UtmSources + `</utm_sources>
      			</urn:setRegistration>
			</soapenv:Body>
		</soapenv:Envelope>`,
	))

	soapAction := "urn:setRegistration"
	httpMethod := "POST"

	req, err := http.NewRequest(httpMethod, url, bytes.NewReader(payload))

	if err != nil {
		return err
	}

	req.Header.Set("Content-type", "text/xml")
	req.Header.Set("SOAPAction", soapAction)

	// Fire the request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	_, err = client.Do(req)

	if err != nil {
		return err
	}

	var quests = `[
		{
			"text": "What is your preferred project?",
			"answer": "` + info.Project + `",
			"type": "Multiple Choice",
			"global_id": "15Zv82NZ3H20220505092400"
		}
	]`

	payload = []byte(strings.TrimSpace(`
		<soapenv:Envelope xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:urn="urn:Service2Controllerwsdl">
   			<soapenv:Header/>
   			<soapenv:Body>
				<urn:setQuest soapenv:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
					<accessToken xsi:type="xsd:string">` + ssmParams["SENTINO_ACCESS_TOKEN"] + `</accessToken>
         			<name xsi:type="xsd:string">` + info.Name + `</name>
         			<email xsi:type="xsd:string">` + info.Email + `</email>
         			<mobile xsi:type="xsd:string">` + info.Phone + `</mobile>
         			<address xsi:type="xsd:string"></address>
         			<questions xsi:type="xsd:string">` + quests + `</questions>
         			<project xsi:type="xsd:string">` + os.Getenv("SENTINO_PROJECT_ID") + `</project>
         			<source xsi:type="xsd:string">` + os.Getenv("SENTINO_SOURCE_ID") + `</source>
      			</urn:setQuest>
			</soapenv:Body>
		</soapenv:Envelope>`,
	))

	soapAction = "urn:setQuest"
	httpMethod = "POST"

	req, err = http.NewRequest(httpMethod, url, bytes.NewReader(payload))

	if err != nil {
		return err
	}

	req.Header.Set("Content-type", "text/xml")
	req.Header.Set("SOAPAction", soapAction)

	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	_, err = client.Do(req)

	if err != nil {
		return err
	}

	return nil
}

/**************************************
SENTINO END
 ***************************************/

func logError(text string) {
	loc, _ := time.LoadLocation("Asia/Kuala_Lumpur")
	now := time.Now().In(loc).Format("2006.01.02 15:04:05")
	write := fmt.Sprintf("[%s] ERROR in %s: %s", now, os.Getenv("FUNCTION_NAME"), text)
	bot, _ := tgbotapi.NewBotAPI(ssmParams["TELEGRAM_BOT_APIKEY"])

	s, _ := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	msg := tgbotapi.NewMessage(s, write)
	bot.Send(msg)

}

/**************************************
GOOGLE RECAPTCHA VERIFICATION START
 ***************************************/
/*func responseRecaptcha(token string) error {
	resp, err := http.Post("https://www.google.com/recaptcha/api/siteverify?secret="+os.Getenv("GOOGLE_RECAPTCHA_SECRET")+"&response="+token, "application/json", nil)
	if err != nil{
		return err
	}

	re := new(recaptchaResponse)

	body, _ := ioutil.ReadAll(resp.Body)
	_ = json.Unmarshal(body, re)
	// fmt.Println(re)

	if re.Success == true{
		return nil
	} else {
		return errors.New(re.ErrorCodes[0])
	}

}*/
/**************************************
GOOGLE RECAPTCHA VERIFICATION END
 ***************************************/

/*
*************************************
MAIN HANDLER START

	**************************************
*/
func StartHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if req.Headers["Content-Type"] != "application/json" &&
		req.Headers["content-type"] != "application/json" {
		return clientError(http.StatusNotAcceptable)
	}

	contact := new(Contact)

	err := json.Unmarshal([]byte(req.Body), contact)

	testing, _ := strconv.ParseBool(os.Getenv("TEST"))
	if !testing {
		os.Setenv("FUNCTION_NAME", "prod-gc-products")
		os.Setenv("SENTINO_PROJECT_ID", "dSawhU9K20221007151041")
		os.Setenv("SENTINO_SOURCE_ID", "450mG0qaN8w20210513180301")

	}

	if err != nil {
		fmt.Println(err.Error())
		return clientError(http.StatusUnprocessableEntity)
	}

	if contact.Email == "" || contact.Name == "" || contact.Phone == "" || contact.Project == "" {
		return clientError(http.StatusBadRequest)
	}

	/*err = responseRecaptcha(contact.RecaptchaToken)
	if err != nil {
		return clientError(http.StatusBadRequest)
	}*/

	sess, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		logError("config Error: " + err.Error())
	}

	ssmsvc := ssm.NewFromConfig(sess)
	param, err := ssmsvc.GetParameters(context.TODO(), &ssm.GetParametersInput{
		Names: []string{
			"MAILGUN_API_KEY",
			"SENTINO_ACCESS_TOKEN",
			"TELEGRAM_BOT_API_KEY",
		},
		WithDecryption: aws.Bool(false),
	})
	if err != nil {
		logError("ssm Error: " + err.Error())
	}

	ssmParams = make(map[string]string)
	values := param.Parameters
	for _, v := range values {
		ssmParams[*v.Name] = *v.Value
	}

	// Send Out Mails
	clientHtml := `<!doctype html>
	<html xmlns="http://www.w3.org/1999/xhtml" xmlns:v="urn:schemas-microsoft-com:vml" xmlns:o="urn:schemas-microsoft-com:office:office">
	   <head>
		  <!--[if gte mso 15]>
		  <xml>
			 <o:OfficeDocumentSettings>
				<o:AllowPNG/>
				<o:PixelsPerInch>96</o:PixelsPerInch>
			 </o:OfficeDocumentSettings>
		  </xml>
		  <![endif]-->
		  <meta charset="UTF-8">
		  <meta http-equiv="X-UA-Compatible" content="IE=edge">
		  <meta name="viewport" content="width=device-width, initial-scale=1">
		  <title>Gamuda Cove</title>
		  <style type="text/css">
			 p{
			 margin:10px 0;
			 padding:0;
			 }
			 table{
			 border-collapse:collapse;
			 }
			 h1,h2,h3,h4,h5,h6{
			 display:block;
			 margin:0;
			 padding:0;
			 }
			 img,a img{
			 border:0;
			 height:auto;
			 outline:none;
			 text-decoration:none;
			 }
			 body,#bodyTable,#bodyCell{
			 height:100%;
			 margin:0;
			 padding:0;
			 width:100%;
			 }
			 .mcnPreviewText{
			 display:none !important;
			 }
			 #outlook a{
			 padding:0;
			 }
			 img{
			 -ms-interpolation-mode:bicubic;
			 }
			 table{
			 mso-table-lspace:0pt;
			 mso-table-rspace:0pt;
			 }
			 .ReadMsgBody{
			 width:100%;
			 }
			 .ExternalClass{
			 width:100%;
			 }
			 p,a,li,td,blockquote{
			 mso-line-height-rule:exactly;
			 }
			 a[href^=tel],a[href^=sms]{
			 color:inherit;
			 cursor:default;
			 text-decoration:none;
			 }
			 p,a,li,td,body,table,blockquote{
			 -ms-text-size-adjust:100%;
			 -webkit-text-size-adjust:100%;
			 }
			 .ExternalClass,.ExternalClass p,.ExternalClass td,.ExternalClass div,.ExternalClass span,.ExternalClass font{
			 line-height:100%;
			 }
			 a[x-apple-data-detectors]{
			 color:inherit !important;
			 text-decoration:none !important;
			 font-size:inherit !important;
			 font-family:inherit !important;
			 font-weight:inherit !important;
			 line-height:inherit !important;
			 }
			 #bodyCell{
			 padding:10px;
			 }
			 .templateContainer{
			 max-width:600px !important;
			 }
			 a.mcnButton{
			 display:block;
			 }
			 .mcnImage,.mcnRetinaImage{
			 vertical-align:bottom;
			 }
			 .mcnTextContent{
			 word-break:break-word;
			 }
			 .mcnTextContent img{
			 height:auto !important;
			 }
			 .mcnDividerBlock{
			 table-layout:fixed !important;
			 }
			 body,#bodyTable{
			 background-color:#FFFFFF;
			 }
			 #bodyCell{
			 border-top:0;
			 }
			 .templateContainer{
			 border:0;
			 }
			 h1{
			 color:#FFFFFF;
			 font-size:28px;
			 font-style:normal;
			 font-weight:bold;
			 line-height:125%;
			 letter-spacing:0;
			 text-align:center;
			 font-family:Helvetica;
			 }
			 h2{
			 color:#FFFFFF;
			 font-family:Helvetica;
			 font-size:22px;
			 font-style:normal;
			 font-weight:bold;
			 line-height:125%;
			 letter-spacing:normal;
			 text-align:center;
			 }
			 h3{
			 color:#FFFFFF;
			 font-family:Helvetica;
			 font-size:20px;
			 font-style:normal;
			 font-weight:bold;
			 line-height:125%;
			 letter-spacing:normal;
			 text-align:center;
			 }
			 h4{
			 color:#FFFFFF;
			 font-size:18px;
			 font-style:normal;
			 font-weight:normal;
			 line-height:125%;
			 letter-spacing:normal;
			 text-align:center;
			 font-family:Helvetica;
			 }
			 #templatePreheader{
			 background-color:#FFFFFF;
			 border-top:0;
			 border-bottom:0;
			 padding-top:9px;
			 padding-bottom:9px;
			 }
			 #templatePreheader .mcnTextContent,#templatePreheader .mcnTextContent p{
			 color:#656565;
			 font-family:Helvetica;
			 font-size:12px;
			 line-height:150%;
			 text-align:center;
			 }
			 #templatePreheader .mcnTextContent a,#templatePreheader .mcnTextContent p a{
			 color:#656565;
			 font-weight:normal;
			 text-decoration:underline;
			 }
			 #templateHeader{
			 background-color:#FFFFFF;
			 border-top:0;
			 border-bottom:0;
			 padding-top:0;
			 padding-bottom:9px;
			 }
			 #templateHeader .mcnTextContent,#templateHeader .mcnTextContent p{
			 color:#202020;
			 font-size:16px;
			 line-height:150%;
			 text-align:left;
			 }
			 #templateHeader .mcnTextContent a,#templateHeader .mcnTextContent p a{
			 color:#2BAADF;
			 font-weight:normal;
			 text-decoration:underline;
			 }
			 #templateBody{
			 background-color:#fff;
			 border-top:0;
			 border-bottom:0;
			 padding-top:0;
			 padding-bottom:0px;
			 }
			 #templateBody .mcnTextContent,#templateBody .mcnTextContent p{
			 color:#FFFFFF;
			 font-size:16px;
			 line-height:150%;
			 text-align:left;
			 }
			 #templateBody .mcnTextContent a,#templateBody .mcnTextContent p a{
			 color:#FFFFFF;
			 font-weight:normal;
			 text-decoration:underline;
			 }
			 #templateFooter{
			 background-color:#FFFFFF;
			 border-top:0;
			 border-bottom:0;
			 padding-top:9px;
			 padding-bottom:9px;
			 }
			 #templateFooter .mcnTextContent,#templateFooter .mcnTextContent p{
			 color:#656565;
			 font-family:Helvetica;
			 font-size:12px;
			 line-height:150%;
			 text-align:center;
			 }
			 #templateFooter .mcnTextContent a,#templateFooter .mcnTextContent p a{
			 color:#656565;
			 font-weight:normal;
			 text-decoration:underline;
			 }
			 @media only screen and (min-width:768px){
			 .templateContainer{
			 width:600px !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 body,table,td,p,a,li,blockquote{
			 -webkit-text-size-adjust:none !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 body{
			 width:100% !important;
			 min-width:100% !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 #bodyCell{
			 padding-top:10px !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 .mcnRetinaImage{
			 max-width:100% !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 .mcnImage{
			 width:100% !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 .mcnCartContainer,.mcnCaptionTopContent,.mcnRecContentContainer,.mcnCaptionBottomContent,.mcnTextContentContainer,.mcnBoxedTextContentContainer,.mcnImageGroupContentContainer,.mcnCaptionLeftTextContentContainer,.mcnCaptionRightTextContentContainer,.mcnCaptionLeftImageContentContainer,.mcnCaptionRightImageContentContainer,.mcnImageCardLeftTextContentContainer,.mcnImageCardRightTextContentContainer,.mcnImageCardLeftImageContentContainer,.mcnImageCardRightImageContentContainer{
			 max-width:100% !important;
			 width:100% !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 .mcnBoxedTextContentContainer{
			 min-width:100% !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 .mcnImageGroupContent{
			 padding:9px !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 .mcnCaptionLeftContentOuter .mcnTextContent,.mcnCaptionRightContentOuter .mcnTextContent{
			 padding-top:9px !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 .mcnImageCardTopImageContent,.mcnCaptionBottomContent:last-child .mcnCaptionBottomImageContent,.mcnCaptionBlockInner .mcnCaptionTopContent:last-child .mcnTextContent{
			 padding-top:18px !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 .mcnImageCardBottomImageContent{
			 padding-bottom:9px !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 .mcnImageGroupBlockInner{
			 padding-top:0 !important;
			 padding-bottom:0 !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 .mcnImageGroupBlockOuter{
			 padding-top:9px !important;
			 padding-bottom:9px !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 .mcnTextContent,.mcnBoxedTextContentColumn{
			 padding-right:18px !important;
			 padding-left:18px !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 .mcnImageCardLeftImageContent,.mcnImageCardRightImageContent{
			 padding-right:18px !important;
			 padding-bottom:0 !important;
			 padding-left:18px !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 .mcpreview-image-uploader{
			 display:none !important;
			 width:100% !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 h1{
			 font-size:22px !important;
			 line-height:125% !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 h2{
			 font-size:20px !important;
			 line-height:125% !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 h3{
			 font-size:18px !important;
			 line-height:125% !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 h4{
			 font-size:16px !important;
			 line-height:150% !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 .mcnBoxedTextContentContainer .mcnTextContent,.mcnBoxedTextContentContainer .mcnTextContent p{
			 font-size:14px !important;
			 line-height:150% !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 #templatePreheader{
			 display:block !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 #templatePreheader .mcnTextContent,#templatePreheader .mcnTextContent p{
			 font-size:14px !important;
			 line-height:150% !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 #templateHeader .mcnTextContent,#templateHeader .mcnTextContent p{
			 font-size:16px !important;
			 line-height:150% !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 #templateBody .mcnTextContent,#templateBody .mcnTextContent p{
			 font-size:16px !important;
			 line-height:150% !important;
			 }
			 }	@media only screen and (max-width: 480px){
			 #templateFooter .mcnTextContent,#templateFooter .mcnTextContent p{
			 font-size:14px !important;
			 line-height:150% !important;
			 }
			 }
		  </style>
	   </head>
	   <body>
		  <!--*|IF:MC_PREVIEW_TEXT|*-->
		  <!--[if !gte mso 9]><!----><span class="mcnPreviewText" style="display:none; font-size:0px; line-height:0px; max-height:0px; max-width:0px; opacity:0; overflow:hidden; visibility:hidden; mso-hide:all;">Our Relationship Associate will be in touch with you soon.</span><!--<![endif]-->
		  <!--*|END:IF|*-->
		  <center>
			 <table align="center" border="0" cellpadding="0" cellspacing="0" height="100%" width="100%" id="bodyTable">
				<tr>
				   <td align="center" valign="top" id="bodyCell">
					  <!-- BEGIN TEMPLATE // -->
					  <!--[if gte mso 9]>
					  <table align="center" border="0" cellspacing="0" cellpadding="0" width="600" style="width:600px;">
						 <tr>
							<td align="center" valign="top" width="600" style="width:600px;">
							   <![endif]-->
							   <table border="0" cellpadding="0" cellspacing="0" width="100%" class="templateContainer">
								  <tr>
									 <td valign="top" id="templateHeader"></td>
								  </tr>
								  <tr>
									 <td valign="top" id="templateBody">
										<table border="0" cellpadding="0" cellspacing="0" width="100%" class="mcnImageBlock" style="min-width:100%;">
										   <tbody class="mcnImageBlockOuter">
											  <tr>
												 <td valign="top" style="padding:0px" class="mcnImageBlockInner">
													<table align="left" width="100%" border="0" cellpadding="0" cellspacing="0" class="mcnImageContentContainer" style="min-width:100%;">
													   <tbody>
														  <tr>
															 <td class="mcnImageContent" valign="top" style="padding-right: 0px; padding-left: 0px; padding-top: 0; padding-bottom: 0; text-align:center;">
																<a href="https://gamudacove.com.my/products" target="_blank">
																<img align="center" alt="" src="https://development.gamudaland.com.my/gamudacove/img/edm/header.png" width="600" style="max-width:600px; padding-bottom: 0; display: inline !important; vertical-align: bottom;" class="mcnImage">
																</a>
															 </td>
														  </tr>
													   </tbody>
													</table>
												 </td>
											  </tr>
										   </tbody>
										</table>
										<table border="0" cellpadding="0" cellspacing="0" width="100%" class="mcnDividerBlock" style="min-width:100%;">
										   <tbody class="mcnDividerBlockOuter">
											  <tr>
												 <td class="mcnDividerBlockInner" style="min-width:100%; padding:18px;">
													<table class="mcnDividerContent" border="0" cellpadding="0" cellspacing="0" width="100%" style="min-width: 100%;border-top: 2px none #EAEAEA;">
													   <tbody>
														  <tr>
															 <td>
																<span></span>
															 </td>
														  </tr>
													   </tbody>
													</table>
													<!--            
													   <td class="mcnDividerBlockInner" style="padding: 18px;">
													   <hr class="mcnDividerContent" style="border-bottom-color:none; border-left-color:none; border-right-color:none; border-bottom-width:0; border-left-width:0; border-right-width:0; margin-top:0; margin-right:0; margin-bottom:0; margin-left:0;" />
													   -->
												 </td>
											  </tr>
										   </tbody>
										</table>
										<table border="0" cellpadding="0" cellspacing="0" width="100%" class="mcnDividerBlock" style="min-width:100%;">
										   <tbody class="mcnDividerBlockOuter">
											  <tr>
												 <td class="mcnDividerBlockInner" style="min-width:100%; padding:18px;">
													<table class="mcnDividerContent" border="0" cellpadding="0" cellspacing="0" width="100%" style="min-width: 100%;border-top: 2px none #EAEAEA;">
													   <tbody>
														  <tr>
															 <td>
																<span></span>
															 </td>
														  </tr>
													   </tbody>
													</table>
													<!--            
													   <td class="mcnDividerBlockInner" style="padding: 18px;">
													   <hr class="mcnDividerContent" style="border-bottom-color:none; border-left-color:none; border-right-color:none; border-bottom-width:0; border-left-width:0; border-right-width:0; margin-top:0; margin-right:0; margin-bottom:0; margin-left:0;" />
													   -->
												 </td>
											  </tr>
										   </tbody>
										</table>
										<table border="0" cellpadding="0" cellspacing="0" width="100%" class="mcnImageBlock" style="min-width:100%;">
										   <tbody class="mcnImageBlockOuter">
											  <tr>
												 <td valign="top" style="padding:9px" class="mcnImageBlockInner">
													<table align="left" width="100%" border="0" cellpadding="0" cellspacing="0" class="mcnImageContentContainer" style="min-width:100%;">
													   <tbody>
														  <tr>
															 <td class="mcnImageContent" valign="top" style="padding-right: 9px; padding-left: 9px; padding-top: 0; padding-bottom: 0; text-align:center;">
																<img align="center" alt="" src="https://development.gamudaland.com.my/gamudacove/img/edm/icon.png" width="102" style="max-width:102px; padding-bottom: 0; display: inline !important; vertical-align: bottom;" class="mcnImage">
															 </td>
														  </tr>
													   </tbody>
													</table>
												 </td>
											  </tr>
										   </tbody>
										</table>
										<table border="0" cellpadding="0" cellspacing="0" width="100%" class="mcnTextBlock" style="min-width:100%;">
										   <tbody class="mcnTextBlockOuter">
											  <tr>
												 <td valign="top" class="mcnTextBlockInner" style="padding-top:9px;">
													<!--[if mso]>
													<table align="left" border="0" cellspacing="0" cellpadding="0" width="100%" style="width:100%;">
													   <tr>
														  <![endif]-->
														  <!--[if mso]>
														  <td valign="top" width="600" style="width:600px;">
															 <![endif]-->
															 <table align="left" border="0" cellpadding="0" cellspacing="0" style="max-width:100%; min-width:100%;" width="100%" class="mcnTextContentContainer">
																<tbody>
																   <tr>
																	  <td valign="top" class="mcnTextContent" style="padding-top:0; padding-right:18px; padding-bottom:9px; padding-left:18px;">
																		 &nbsp;
																		 <h1><span style="color:#3f2a27">We have received your registration!</span></h1>
																		 &nbsp;
																		 <h4><span style="color:#292929">Our Relationship Associate will be in touch with you soon.</span></h4>
																	  </td>
																   </tr>
																</tbody>
															 </table>
															 <!--[if mso]>
														  </td>
														  <![endif]-->
														  <!--[if mso]>
													   </tr>
													</table>
													<![endif]-->
												 </td>
											  </tr>
										   </tbody>
										</table>
										<table border="0" cellpadding="0" cellspacing="0" width="100%" class="mcnDividerBlock" style="min-width:100%;">
										   <tbody class="mcnDividerBlockOuter">
											  <tr>
												 <td class="mcnDividerBlockInner" style="min-width:100%; padding:18px;">
													<table class="mcnDividerContent" border="0" cellpadding="0" cellspacing="0" width="100%" style="min-width: 100%;border-top: 2px none #EAEAEA;">
													   <tbody>
														  <tr>
															 <td>
																<span></span>
															 </td>
														  </tr>
													   </tbody>
													</table>
												 </td>
											  </tr>
										   </tbody>
										</table>
										<table border="0" cellpadding="0" cellspacing="0" width="100%" class="mcnDividerBlock" style="min-width:100%;">
										   <tbody class="mcnDividerBlockOuter">
											  <tr>
												 <td class="mcnDividerBlockInner" style="min-width:100%; padding:18px;">
													<table class="mcnDividerContent" border="0" cellpadding="0" cellspacing="0" width="100%" style="min-width: 100%;border-top: 2px none #EAEAEA;">
													   <tbody>
														  <tr>
															 <td>
																<span></span>
															 </td>
														  </tr>
													   </tbody>
													</table>
												 </td>
											  </tr>
										   </tbody>
										</table>
									 </td>
								  </tr>
								  <tr>
									 <td align="center" valign="top">
										<a href="https://gamudacove.com.my/" target="_blank"> <img src="https://development.gamudaland.com.my/gamudacove/img/edm/footer.png" width="100%" height="auto" style="display:block; margin:0; padding:0;">
										</a>
									 </td>
								  </tr>
								  <tr>
									 <td valign="top" id="templateFooter">
										<table border="0" cellpadding="0" cellspacing="0" width="100%" class="mcnTextBlock" style="min-width:100%;">
										   <tbody class="mcnTextBlockOuter">
											  <tr>
												 <td valign="top" class="mcnTextBlockInner" style="padding-top:9px;">
													<!--[if mso]>
													<table align="left" border="0" cellspacing="0" cellpadding="0" width="100%" style="width:100%;">
													   <tr>
														  <![endif]-->
														  <!--[if mso]>
														  <td valign="top" width="600" style="width:600px;">
															 <![endif]-->
															 <table align="left" border="0" cellpadding="0" cellspacing="0" style="max-width:100%; min-width:100%;" width="100%" class="mcnTextContentContainer">
																<tbody>
																   <tr>
																	  <td valign="top" class="mcnTextContent" style="padding-top:0; padding-right:18px; padding-bottom:9px; padding-left:18px;">
																		 <em>Copyright © Gamuda Cove. All rights reserved.</em>
																	  </td>
																   </tr>
																   <tr>
																	  <td>
																		 <p style="font-size: 12px;font-family: OpenSans, sans-serif;text-align: center;color: #000;padding: 0 25px">This inbox is not monitored, please do not reply to this email.</p>
																	  </td>
																   </tr>
																</tbody>
															 </table>
															 <!--[if mso]>
														  </td>
														  <![endif]-->
														  <!--[if mso]>
													   </tr>
													</table>
													<![endif]-->
												 </td>
											  </tr>
										   </tbody>
										</table>
									 </td>
								  </tr>
							   </table>
							   <!--[if gte mso 9]>
							</td>
						 </tr>
					  </table>
					  <![endif]-->
					  <!-- // END TEMPLATE -->
				   </td>
				</tr>
			 </table>
		  </center>
	   </body>
	</html>`

	clientPlain := ""

	_, err = SendSimpleMail(contact.Email, "Hi "+contact.Name+", thank you for your interest in Gamuda Cove!", clientHtml, clientPlain, "Gamuda Cove | Gamuda Land", "noreply@gamudaland.com.my")

	if err != nil {
		log.Printf("MailError: %v", err)
		// return serverError(err)
	}

	// Save Details to DynamoDB
	err = SaveToDynamoDb(contact)

	if err != nil {
		log.Printf("DynamoError: %v", err)
		return serverError(err)
	}

	// Save to Sentino
	err = SaveToSentino(contact)

	if err != nil {
		log.Printf("SentinoError: %v", err)
		return serverError(err)
	}

	loc, _ := time.LoadLocation("Asia/Kuala_Lumpur")
	t := time.Now().In(loc)

	layout := "2006-01-02 15:04:05"

	adminHtml := `<html>
		` + contact.Name + ` has registered in Gamuda Cove Product Landing Page at ` + t.Format(layout) + ` Details as follow:
		<br>
		<h3>Name: ` + contact.Name + `</h3>
		<h3>Phone Number: ` + contact.Phone + `</h3>
		<h3>Email: ` + contact.Email + `</h3>
		<h3>Project: ` + contact.Project + `</h3>
	</html>`

	adminPlain := ""

	_, err = SendSimpleMail(os.Getenv("MG_ADMIN_MAIL"), os.Getenv("MG_ADMIN_SUBJECT"), adminHtml, adminPlain, "Gamuda Cove | Gamuda Land", "noreply@gamudaland.com.my")

	if err != nil {
		log.Printf("AdminMailError: %v", err)
		// return serverError(err)
	}

	body, _ := json.Marshal(contact)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(body),
		Headers:    map[string]string{"Content-Type": "application/json", "Access-Control-Allow-Origin": "*"},
	}, nil
}

/**************************************
MAIN HANDLER END
 ***************************************/

/*
*
ERROR STUFF
*/
var errorLogger = log.New(os.Stderr, "ERROR", log.Llongfile)

func serverError(err error) (events.APIGatewayProxyResponse, error) {
	logError(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    map[string]string{"Access-Control-Allow-Origin": "*"},
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(StartHandler)
}
