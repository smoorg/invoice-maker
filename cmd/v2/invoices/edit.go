package invoices

import (
	labelinput "invoice-maker/cmd/v2/controls/label_input"
	singleselect "invoice-maker/cmd/v2/controls/single_select"
	"strings"

	"invoice-maker/pkg/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	FocusNone = iota
	FocusInvoiceNo
	FocusInvoiceDate
	FocusDeliveryDate
	FocusDueDate
	FocusReceiver
	FocusPaymentType
	// FocusItems
)

const FocusNumOfItems = 4

type InvoiceEditModel struct {
	focus     int
	invoice   *config.Invoice
	receivers []config.Company
	issuers   []config.Issuer

	InvoiceNo    labelinput.Model
	InvoiceDate  labelinput.Model
	DeliveryDate labelinput.Model
	DueDate      labelinput.Model
	Receiver     singleselect.Model
	PaymentType  singleselect.Model
	// Items            singleselect.Model
}

func (m *InvoiceEditModel) SetInvoice(v *config.Invoice) {
	m.invoice = v
	m.InvoiceNo.Input.SetValue(v.InvoiceNo)
	m.InvoiceDate.Input.SetValue(v.InvoiceDate)
	m.DeliveryDate.Input.SetValue(v.DeliveryDate)
	m.DueDate.Input.SetValue(v.DueDate)
	m.Receiver.SetValue(v.Receiver.Name)
	m.PaymentType.SetValue(v.PaymentType)
	// m.Issuer.SetValue(v.Issuer.Name)
	// m.Items.Input.SetValue(v.InvoiceNo)
}

func NewEditModel() InvoiceEditModel {
	m := InvoiceEditModel{}

	m.InvoiceNo = labelinput.New("Invoice No.", 20)
	m.InvoiceDate = labelinput.New("Invoice Date", 20)
	m.DeliveryDate = labelinput.New("DeliveryDate", 20)
	m.DueDate = labelinput.New("Due Date", 20)

	cfg, err := config.GetConfig(nil)
	if err != nil {
		panic(err)
	}
	recs := []singleselect.Item{}
	for _, v := range cfg.Receivers {
		recs = append(recs, singleselect.Item{Label: v.Name, Value: v.Name})
	}
	m.Receiver = singleselect.New("Receiver", 20, recs)

	paymentTypes := []singleselect.Item{
		{Label: "Cash", Value: "Cash"},
		{Label: "Transfer", Value: "Transfer"},
	}
	m.PaymentType = singleselect.New("Payment Type", 20, paymentTypes)

	return m
}
func (m InvoiceEditModel) Init() {
}

func (m InvoiceEditModel) Update(msg tea.Msg) (InvoiceEditModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case tea.KeyTab.String():
			cmd = m.CycleFocus(1)
			cmds = append(cmds, cmd)
			return m, tea.Batch(cmds...)
		case tea.KeyShiftTab.String():
			cmd = m.CycleFocus(-1)
			cmds = append(cmds, cmd)
		}
	}

	switch m.focus {
	case FocusNone:
	case FocusInvoiceNo:
		m.InvoiceNo, cmd = m.InvoiceNo.Update(msg)
		cmds = append(cmds, cmd)
	case FocusInvoiceDate:
		m.InvoiceDate, cmd = m.InvoiceDate.Update(msg)
		cmds = append(cmds, cmd)
	case FocusDeliveryDate:
		m.DeliveryDate, cmd = m.DeliveryDate.Update(msg)
		cmds = append(cmds, cmd)
	case FocusDueDate:
		m.DueDate, cmd = m.DueDate.Update(msg)
		cmds = append(cmds, cmd)
	case FocusReceiver:
		m.Receiver, cmd = m.Receiver.Update(msg)
		cmds = append(cmds, cmd)
	case FocusPaymentType:
		m.PaymentType, cmd = m.PaymentType.Update(msg)
		cmds = append(cmds, cmd)
		// case FocusItems:
		// 	m.Items, cmd = m.Items.Update(msg)
		// 	cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *InvoiceEditModel) Blur() {
	m.InvoiceNo.Blur()
	m.InvoiceDate.Blur()
	m.DeliveryDate.Blur()
	m.DueDate.Blur()
	m.Receiver.Blur()
	m.PaymentType.Blur()
	// m.Items.Blur()
}

func (m *InvoiceEditModel) CycleFocus(increment int) tea.Cmd {
	m.Blur()
	m.focus = m.focus + increment
	switch m.focus {
	case FocusInvoiceNo:
		return m.InvoiceNo.Focus()
	case FocusInvoiceDate:
		return m.InvoiceDate.Focus()
	case FocusDeliveryDate:
		return m.DeliveryDate.Focus()
	case FocusDueDate:
		return m.DueDate.Focus()
	case FocusReceiver:
		m.Receiver.Focus()
	case FocusPaymentType:
		m.PaymentType.Focus()
	default:
		if increment < 0 {
			// we're going back, land on the last
			m.focus = FocusPaymentType
			m.PaymentType.Focus()
		} else {
			m.focus = 0
		}
	}

	return nil
}

func (m InvoiceEditModel) View() string {
	s := strings.Builder{}
	formItems := []string{
		m.InvoiceNo.View(),
		m.InvoiceDate.View(),
		m.DeliveryDate.View(),
		m.DueDate.View(),
		m.Receiver.View(),
		m.PaymentType.View(),
		//m.Items.View(),
	}

	for i, v := range formItems {
		s.WriteString(v)

		// avoid last
		if i < len(formItems)-1 {
			s.WriteString("\n")
			s.WriteString("\n")
		}
	}

	content := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).Padding(1).Render(s.String())

	return content
}
