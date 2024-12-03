package smsmessage

import "messageapi.e-vrit.co.il/services/smsmessage/models"

type ISms interface {
	SendSms(models.SmsModel) (bool, error)
}
