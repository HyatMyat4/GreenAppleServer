package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PaymentMethods struct {
	Id                  string `json:"id"`
	Payment_method_name string `json:"payment_method_name" validate:"required"`
}

type VoidReasons struct {
	Id          string `json:"id"`
	Reason_name string `json:"void_reason_name" validate:"required"`
}

type Description_translate struct{}

type Comapny struct {
	Id                 primitive.ObjectID      `bson:"_id"`
	Company_id         string                  `json:"company_id" validate:"required"`
	Company_name       string                  `json:"name"`
	Description        []Description_translate `json:"description"`
	Image_url          string                  `json:"image_url"`
	Logo_url           string                  `json:"logo_url"`
	Company_color      string                  `json:"company_color"`
	Currency           string                  `json:"currency" validate:"eq=USD|eq=THB|eq=MMK|eq=SGD|eq=CNY|eq=JPY|eq=KWR|eq=INR"`
	Stripe_customer_id string                  `json:"strupe_coustomer_id" validate:"required"`
	Subscription_id    string                  `json:"subscription_id"`
	Selected_languages string                  `json:"selected_languages" validate:"eq=en|eq=th|eq=my|eq=zh|eq=ja|eq=ko|eq=hi"`
	Tax_amount         string                  `json:"tax_amount"`
	Services_charge    string                  `json:"services_charge"`
	Payment_methods    []PaymentMethods        `json:"payment_methods"`
	Void_reasons       []VoidReasons           `json:"void_reasons"`
}
