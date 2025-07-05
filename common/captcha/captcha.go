package captcha

import "github.com/mojocn/base64Captcha"

type Captcha struct {
	captcha *base64Captcha.Captcha
}

// Initialize captcha
func NewCaptcha() *Captcha {

	driver := base64Captcha.NewDriverDigit(40, 100, 4, 0.7, 1)

	return &Captcha{
		captcha: base64Captcha.NewCaptcha(driver, &RedisStore{}),
	}
}

// Generate captcha
// uuid, base64, answer
func (c *Captcha) Generate() (string, string) {

	id, b64s, _, err := c.captcha.Generate()
	if err != nil {
		return "", ""
	}

	return id, b64s
}

// Verify captcha
func (c *Captcha) Verify(id, answer string) bool {
	return c.captcha.Verify(id, answer, true)
}
