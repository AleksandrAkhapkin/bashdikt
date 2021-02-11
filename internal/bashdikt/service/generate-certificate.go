package service

import (
	"fmt"
	"github.com/AleksandrAkhapkin/bashdikt/internal/bashdikt/types"
	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"strings"
)

func (s *Service) makeAndSendCertificate(cert *types.MakeCertificate) error {

	const X = 1779 //размеры файла
	const Y = 2480
	fontsSize := 70.0
	nameFiles := fmt.Sprintf("%s%s%s", cert.LastName, cert.FirstName, cert.MiddleName)
	nameFiles = strings.ReplaceAll(nameFiles, " ", "")
	pathForSave := fmt.Sprintf("%s/%s.png", cert.PathForSave, nameFiles)

	//выбираем для кого сертификат
	if strings.Contains(cert.PathForCert, "student") {

		//для студента
		im, err := gg.LoadImage(cert.PathForCert)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("err with gg.LoadImage in MakeAndSendCertificate FOR USER_ID = %d", cert.UserID))
		}

		dc := gg.NewContext(X, Y)
		dc.Clear()
		dc.SetRGB(0, 0, 0)
		if err := dc.LoadFontFace(cert.PathForFonts, fontsSize); err != nil {
			return errors.Wrap(err, fmt.Sprintf("err with dc.LoadFontFace in MakeAndSendCertificate FOR USER_ID = %d", cert.UserID))
		}

		dc.DrawRoundedRectangle(0, 0, 512, 512, 0)
		dc.DrawImage(im, 0, 0)
		dc.DrawStringAnchored(cert.LastName, X/1.7, 1080, 0.5, 0.5)
		dc.DrawStringAnchored(cert.FirstName, X/1.7, 1160, 0.5, 0.5)
		dc.DrawStringAnchored(cert.MiddleName, X/1.7, 1240, 0.5, 0.5)

		dc.Clip()
		if err = dc.SavePNG(pathForSave); err != nil {
			return errors.Wrap(err, fmt.Sprintf("err with dc.SavePNG in MakeAndSendCertificate FOR USER_ID = %d", cert.UserID))
		}

	} else {

		//для учителя
		im, err := gg.LoadImage(cert.PathForCert)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("err with gg.LoadImage in MakeAndSendCertificate FOR USER_ID = %d", cert.UserID))
		}

		dc := gg.NewContext(X, Y)
		dc.Clear()
		dc.SetRGB(0, 0, 0)
		if err := dc.LoadFontFace(cert.PathForFonts, fontsSize); err != nil {
			return errors.Wrap(err, fmt.Sprintf("err with dc.LoadFontFace in MakeAndSendCertificate FOR USER_ID = %d", cert.UserID))
		}
		dc.DrawRoundedRectangle(0, 0, 512, 512, 0)
		dc.DrawImage(im, 0, 0)
		dc.DrawStringAnchored(cert.LastName, X/1.7, 840, 0.5, 0.5)
		dc.DrawStringAnchored(cert.FirstName, X/1.7, 920, 0.5, 0.5)
		dc.DrawStringAnchored(cert.MiddleName, X/1.7, 1000, 0.5, 0.5)

		dc.Clip()
		if err = dc.SavePNG(pathForSave); err != nil {
			return errors.Wrap(err, fmt.Sprintf("err with dc.SavePNG in MakeAndSendCertificate FOR USER_ID = %d", cert.UserID))
		}
	}

	if err := s.email.SendCert(pathForSave, cert.EmailTo); err != nil {
		return errors.Wrap(err, fmt.Sprintf("err with SendCert in MakeAndSendCertificate FOR USER_ID = %d", cert.UserID))
	}

	return nil
}
