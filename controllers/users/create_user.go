package users

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"green-apple-server/database"
	"green-apple-server/helper"
	"green-apple-server/mail"
	"green-apple-server/models"
	"green-apple-server/token"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RequestBody struct {
	User_name          string    `json:"user_name" validate:"required,min=2,max=80"`
	Password           string    `json:"password" validate:"required,min=8,max=30"`
	Email              string    `json:"email" validate:"email,required"`
	Phone              string    `json:"phone" validate:"required,min=5,max=30"`
	Role               string    `json:"role" validate:"required,eq=admin|eq=user"`
	Created_at         time.Time `json:"created_at"`
	Updated_at         time.Time `json:"updated_at"`
	Company_name       string    `json:"name"`
	Currency           string    `json:"currency" validate:"eq=USD|eq=THB|eq=MMK|eq=SGD|eq=CNY|eq=JPY|eq=KWR|eq=INR"`
	Selected_languages string    `json:"selected_languages" validate:"eq=en|eq=th|eq=my|eq=zh|eq=ja|eq=ko|eq=hi"`
}

type Create_user_response struct {
	Id  string `json:"_id"`
	Opt string `json:"opt"`
}

var validate = validator.New()

var payments = []models.PaymentMethods{
	{Id: xid.New().String(), Payment_method_name: "Cash"},
	{Id: xid.New().String(), Payment_method_name: "Credit Card"},
	{Id: xid.New().String(), Payment_method_name: "QR Code"},
}

var void_reason = []models.VoidReasons{
	{Id: xid.New().String(), Reason_name: "Customer Error"},
	{Id: xid.New().String(), Reason_name: "Staff Error"},
	{Id: xid.New().String(), Reason_name: "Kitchen Error"},
	{Id: xid.New().String(), Reason_name: "Out of Stock"},
}

var wg = &sync.WaitGroup{}

var mut = &sync.Mutex{}

var access_Token string

func Generatetoken(reqbody RequestBody, mut *sync.Mutex, user models.User) {
	wg.Add(1)
	mut.Lock()
	access_Token = token.GenerateAllTokens(reqbody.Email, reqbody.User_name, user.Id.Hex())
	mut.Unlock()
	wg.Done()
}

func CheckEmail(mut *sync.Mutex, ctx context.Context, reqbody RequestBody, req *gin.Context) {
	wg.Add(1)
	mut.Lock()
	email_count, err := database.UsersCollection.CountDocuments(ctx, bson.M{"email": reqbody.Email})

	if err != nil {
		req.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
		log.Panic(err)
		return
	}

	if email_count > 0 {
		req.JSON(http.StatusOK, gin.H{"message": "email already exists!"})
		return
	}
	mut.Unlock()
	wg.Done()
}

func InsertUser(mut *sync.Mutex, insertuserCh chan *mongo.InsertOneResult, ctx context.Context, user models.User, req *gin.Context) {
	wg.Add(1)
	mut.Lock()
	result, insert_error := database.UsersCollection.InsertOne(ctx, user)
	if insert_error != nil {
		req.JSON(http.StatusInternalServerError, gin.H{"message": insert_error.Error()})
		return
	}
	insertuserCh <- result
	mut.Unlock()
	wg.Done()
}

func Create_user() gin.HandlerFunc {
	return func(req *gin.Context) {

		var insertuserChannel = make(chan *mongo.InsertOneResult)

		var optChannel = make(chan string)

		var emailSenderChannel = make(chan error)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var reqbody RequestBody

		var user models.User

		var company models.Comapny

		if err := req.BindJSON(&reqbody); err != nil {
			defer cancel()
			req.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		user.Id = primitive.NewObjectID()
		user.User_name = reqbody.User_name
		user.UserId = user.Id
		user.Email = reqbody.Email

		go Generatetoken(reqbody, mut, user)

		validationErr := validate.Struct(reqbody)

		if validationErr != nil {
			req.JSON(http.StatusBadRequest, gin.H{"message": validationErr.Error()})
			defer cancel()
			return
		}

		go CheckEmail(mut, ctx, reqbody, req)
		defer cancel()

		phonenumber_count, _ := database.UsersCollection.CountDocuments(ctx, bson.M{"phone": reqbody.Phone})
		defer cancel()

		if phonenumber_count > 0 {
			req.JSON(http.StatusOK, gin.H{"message": "phone number already exists!"})
			return
		}

		go EncodeToString(mut, optChannel, 6)

		OPT := <-optChannel

		user.Password = helper.HashPassword(reqbody.Password)
		user.Phone = reqbody.Phone
		user.Role = reqbody.Role
		user.Token = access_Token
		user.Pin = OPT
		user.Email_verified = false
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		go InsertUser(mut, insertuserChannel, ctx, user, req)

		// result, insert_error := database.UsersCollection.InsertOne(ctx, user)
		// defer cancel()

		// if insert_error != nil {
		// 	req.JSON(http.StatusInternalServerError, gin.H{"message": insert_error.Error()})
		// 	return
		// }

		go VerifyEmailSender(emailSenderChannel, OPT)

		_err := <-emailSenderChannel

		if _err != nil {
			req.JSON(http.StatusInternalServerError, gin.H{"message": "email failed to send"})
			return
		}

		result := <-insertuserChannel

		user_id := result.InsertedID.(primitive.ObjectID).Hex()
		company.Id = primitive.NewObjectID()
		company.Company_id = user_id
		company.Currency = reqbody.Currency
		company.Selected_languages = reqbody.Selected_languages
		company.Stripe_customer_id = OPT
		company.Subscription_id = OPT
		company.Payment_methods = payments
		company.Void_reasons = void_reason

		_, insinsert_error := database.CompanyCollection.InsertOne(ctx, company)
		defer cancel()

		if insinsert_error != nil {
			req.JSON(http.StatusInternalServerError, gin.H{"message": insinsert_error.Error()})
			return
		}

		// Set cookie with max age 1 year (365 days)
		maxAge := 3600 * 24 * 365

		req.SetSameSite(http.SameSiteNoneMode)
		req.SetCookie("jwt_token", access_Token, maxAge, "/", "localhost", true, true)
		req.Header("Access-Control-Allow-Credentials", "true")

		var response Create_user_response
		response.Id = user_id
		response.Opt = OPT

		wg.Wait()

		req.JSON(http.StatusOK, gin.H{"data": response})
	}
}

func VerifyEmailSender(emailSenderChannel chan error, OPT string) {
	wg.Add(1)
	EMAIL_SENDER_NAME := os.Getenv("EMAIL_SENDER_NAME")
	EMAIL_SENDER_ADDRESS := os.Getenv("EMAIL_SENDER_ADDRESS")
	EMAIL_SENDER_PASSWORD := os.Getenv("EMAIL_SENDER_PASSWORD")
	sender := mail.NewGmailSender(EMAIL_SENDER_NAME, EMAIL_SENDER_ADDRESS, EMAIL_SENDER_PASSWORD)

	subject := "Email Verification"

	content, _err := template.New("").Parse(`
	<div style="font-family: Helvetica,Arial,sans-serif;width:100%;overflow:auto;line-height:2">
	<div style="margin:10px auto;width:80%;padding:20px 0">
		<div style="border-bottom:1px solid #eee">
		<a href="" style="font-family: Cedarville Cursive, cursive;font-size:1.4em;color: #79C523;text-decoration:none;font-weight:600">Green Apple</a>
		</div>
		<p style="font-size:1.1em">Hi,</p>
		<p>Use the following OTP to complete your sign up process. OTP is valid for 5 minutes</p>
		<h2 style="background: #79C523;margin: 0 auto;width: max-content;padding: 0 10px;color: #fff;border-radius: 4px;">{{.OPT}}</h2>
		<p style="font-size:0.9em;">Regards,<p style="font-family: Cedarville Cursive, cursive;" > Green Apple</p></p>
		<hr style="border:none;border-top:1px solid #eee" />
		<div style="float:right;padding:8px 0;color:#aaa;font-size:0.8em;line-height:1;font-weight:300">
		<p>Thanks For Using Our App</p>
		<p>Green Apple 2023 Inc.</p>
		<p>Yangon,Myammar</p>
		</div>
	</div>
	</div>
	`)

	if _err != nil {
		fmt.Println(_err)
	}

	var contentBuffer bytes.Buffer
	data := struct {
		OPT string
	}{
		OPT: OPT, // Pass the OPT value to the template
	}

	// Execute the template with the data and store it in the buffer
	_err = content.Execute(&contentBuffer, data)
	if _err != nil {
		fmt.Println(_err)
	}
	to := []string{"rapperlay2584@gmail.com"}

	err := sender.SendEmail(subject, contentBuffer.String(), to, nil, nil)

	emailSenderChannel <- err

	wg.Done()
}

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func EncodeToString(mut *sync.Mutex, optChannel chan string, max int) {
	wg.Add(1)
	mut.Lock()
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	optChannel <- string(b)
	mut.Unlock()
	wg.Done()
}
