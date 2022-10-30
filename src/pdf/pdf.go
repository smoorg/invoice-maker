package pdf

import (
	"context"
	"os"
	"path/filepath"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func PrintInvoice(htmlFileName string, pdfFileName string) error {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	return chromedp.Run(ctx,
		chromedp.Navigate("file://"+htmlFileName),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().WithPrintBackground(false).Do(ctx)
			if err != nil {
				return err
			}
			err = os.WriteFile(filepath.Join(filepath.Dir(htmlFileName), pdfFileName), buf, 0744)
			if err != nil {
				return err
			}

			return nil
		}),
	)
}
