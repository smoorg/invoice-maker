package invoices

import (
	"invoice-maker/cmd/v2/controls/button"
	labelinput "invoice-maker/cmd/v2/controls/label_input"
	singleselect "invoice-maker/cmd/v2/controls/single_select"
	"invoice-maker/pkg"
	"strings"
	"time"

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
	FocusBtnAddItem
	FocusBtnSave
	// FocusItems
)

const FocusNumOfItems = 4

type InvoiceEditModel struct {
	focus             int
	timeLayout        string
	invoice           *config.Invoice
	receivers         []config.Company
	issuers           []config.Issuer

	InputInvoiceNo     labelinput.Model
	InputInvoiceDate   DateInput
	InputDeliveryDate  labelinput.Model
	InputDueDate       labelinput.Model
	SelectReceiverName singleselect.Model
	SelectPaymentType  singleselect.Model

	BtnSave    button.ButtonModel
	BtnAddItem button.ButtonModel
	// Items            singleselect.Model
}

func (m *InvoiceEditModel) SetInvoiceModel(v *config.Invoice) {
	m.invoice = v
	m.InputInvoiceNo.SetValue(v.InvoiceNo)
	m.InputInvoiceDate.SetValue(v.InvoiceDate)
	m.InputDeliveryDate.SetValue(v.DeliveryDate)
	m.InputDueDate.SetValue(v.DueDate)
	m.SelectReceiverName.SetValue(v.Receiver.Name)
	m.SelectPaymentType.SetValue(v.PaymentType)
	// m.Issuer.SetValue(v.Issuer.Name)
	// m.Items.Input.SetValue(v.InvoiceNo)
}

type Option func(*InvoiceEditModel)

func WithLayout(layout string) Option {
	return func(iem *InvoiceEditModel) {
		iem.timeLayout = layout
	}
}

func NewEditModel(opts ...Option) InvoiceEditModel {
	m := InvoiceEditModel{}

	for _, o := range opts {
		o(&m)
	}

	m.InputInvoiceNo = labelinput.New("Invoice No.", 20)
	m.timeLayout = time.DateOnly
	m.InputInvoiceDate = NewDateInput("Invoice Date", 20, time.Time{}, m.timeLayout)
	m.InputDeliveryDate = labelinput.New("DeliveryDate", 20)
	m.InputDueDate = labelinput.New("Due Date", 20)

	cfg, err := config.GetConfig(nil)
	if err != nil {
		panic(err)
	}
	recs := []singleselect.Item{}
	for _, v := range cfg.Receivers {
		recs = append(recs, singleselect.Item{Label: v.Name, Value: v.Name})
	}
	m.SelectReceiverName = singleselect.New("Receiver", 20, recs)

	paymentTypes := []singleselect.Item{
		{Label: "Cash", Value: "Cash"},
		{Label: "Transfer", Value: "Transfer"},
	}
	m.SelectPaymentType = singleselect.New("Payment Type", 20, paymentTypes)

	m.BtnSave = button.NewButton("save", "Save")
	m.BtnAddItem = button.NewButton("add", "Add Item")

	return m
}
func (m InvoiceEditModel) Init() {}

func (m InvoiceEditModel) Update(msg tea.Msg) (InvoiceEditModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case button.ButtonSubmitEvent:
		switch msg.ButtonName {
		case m.BtnSave.Name:
			return m, m.UpdateInvoice()
		case m.BtnAddItem.Name:
			return m, m.AddItem()
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "h":
			if m.focus == 0 {
				return m, pkg.GoInvoicePreview()
			}
		case tea.KeyEsc.String():
			if m.focus == 0 {
				return m, pkg.GoInvoicePreview()
			} else {
				m.Blur()
				m.focus = 0
			}

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
		m.InputInvoiceNo, cmd = m.InputInvoiceNo.Update(msg)
		cmds = append(cmds, cmd)
	case FocusInvoiceDate:
		m.InputInvoiceDate, cmd = m.InputInvoiceDate.Update(msg)
		cmds = append(cmds, cmd)
	case FocusDeliveryDate:
		m.InputDeliveryDate, cmd = m.InputDeliveryDate.Update(msg)
		cmds = append(cmds, cmd)
	case FocusDueDate:
		m.InputDueDate, cmd = m.InputDueDate.Update(msg)
		cmds = append(cmds, cmd)
	case FocusReceiver:
		m.SelectReceiverName, cmd = m.SelectReceiverName.Update(msg)
		cmds = append(cmds, cmd)
	case FocusPaymentType:
		m.SelectPaymentType, cmd = m.SelectPaymentType.Update(msg)
		cmds = append(cmds, cmd)
		// case FocusItems:
		// 	m.Items, cmd = m.Items.Update(msg)
		// 	cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *InvoiceEditModel) AddItem() tea.Cmd {
	m.invoice.Items = append(m.invoice.Items, config.InvoiceItem{})

	return nil
}

func (m *InvoiceEditModel) UpdateInvoice() tea.Cmd {
	m.invoice.InvoiceNo = m.InputInvoiceNo.Value()
	m.invoice.InvoiceDate = m.InputInvoiceDate.Value()
	m.invoice.DeliveryDate = m.InputDeliveryDate.Value()
	m.invoice.DueDate = m.InputDueDate.Value()
	m.invoice.PaymentType = m.SelectPaymentType.Value()

	receiverName := m.SelectReceiverName.Value()
	for _, v := range m.receivers {
		if v.Name == receiverName {
			m.invoice.Receiver = v
		}
	}

	return nil
}

func (m *InvoiceEditModel) Blur() {
	m.InputInvoiceNo.Blur()
	m.InputInvoiceDate.input.Blur()
	m.InputDeliveryDate.Blur()
	m.InputDueDate.Blur()
	m.SelectReceiverName.Blur()
	m.SelectPaymentType.Blur()
	m.BtnAddItem.Blur()
	m.BtnSave.Blur()
	// m.Items.Blur()
}

func (m *InvoiceEditModel) CycleFocus(increment int) tea.Cmd {
	m.Blur()
	m.focus = m.focus + increment
	switch m.focus {
	case FocusInvoiceNo:
		return m.InputInvoiceNo.Focus()
	case FocusInvoiceDate:
		return m.InputInvoiceDate.Focus()
	case FocusDeliveryDate:
		return m.InputDeliveryDate.Focus()
	case FocusDueDate:
		return m.InputDueDate.Focus()
	case FocusReceiver:
		m.SelectReceiverName.Focus()
	case FocusPaymentType:
		m.SelectPaymentType.Focus()
	case FocusBtnAddItem:
		m.BtnAddItem.Focus()
	case FocusBtnSave:
		m.BtnSave.Focus()
	default:
		if increment < 0 {
			// we're going back, land on the last
			m.focus = FocusPaymentType
			m.SelectPaymentType.Focus()
		} else {
			m.focus = 0
		}
	}

	return nil
}

var thickBorderStyle = lipgloss.NewStyle().
	Border(lipgloss.ThickBorder()).
	Padding(1)

var marginStyle = lipgloss.NewStyle().Padding(1)

func (m InvoiceEditModel) View() string {
	s := strings.Builder{}
	formItems := []string{
		m.InputInvoiceNo.View(),
		m.InputInvoiceDate.View(),
		m.InputDeliveryDate.View(),
		m.InputDueDate.View(),
		m.SelectReceiverName.View(),
		m.SelectPaymentType.View(),
		m.BtnAddItem.View(),
		m.BtnSave.View(),
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

	//content := thickBorderStyle.Render(s.String())

	return marginStyle.Render(s.String())
}
