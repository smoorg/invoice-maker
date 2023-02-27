# Invoice Maker

## Introduction 

Program is designed as small CLI tool to help medium and small businesses with tedious task of creating recurring invoices. Can be used to 

- avoid VAT tax calculation (calculates VAT and total amount based on VAT rate provided)
- avoid rounding issues (uses decimals)
- store existing customer's data 
- minimize mistakes in the invoices
- see historical data in single place

It can be used as database of all your partners, or just a simple way to generate invoices from markdown/html format into PDF as most people expects that format.

Keep in mind not everything works as expected. I am working on it in my spare time just for fun.

## How to use it

Makefile is prepared to use it straight out of root folder. Alternatively, all the code is in `src` so you can `cd` to it and just `go build` instead. 


### Build
```sh
make build
```


### Run
```sh
make run
```

### Test

```sh
make test
```

## Getting started

1. For the time being all you need lies in `/templates` folder on the root of this project. Content of it is supposed to land in `~/.config/invoice-maker` and be modified as necessary.
1. After initial run please go to `Config` and specify where your printed invoices should land to i.e. `/home/user/Documents/invoices`. Please also make sure you have write privileges to that folder.

## Keys

For the most part, you can follow Vim keys or arrows as a fallback. 

- `Enter` or `l` confirms
- `ESC` or `h` cancels
- `Arrow` `UP` or `k` moves up
- `Arrow` `DOWN` or `j` moves down
- `Ctrl+D` or `Ctrl+C` to close the app

> Vim keys won't work with forms as we need text input. Alternatively, 
> * Tab or Enter can be used to go to next form control
> * Shift+Tab or ESC can be used to go to previous form control

## Config files


List of cryptic names for the columns in the invoice:

* InvoiceNo		- arbitrary number of your invoice. I tend to use [year]/[month]/[invoiceNo] but it can be anything else.
* InvoiceDate		- date you created that invoice
* DeliveryDate		- date of delivery of the product/service. I usually put month I've been working for this and that company.
* IssuerName		- name of the issuer company, most likely your company
* ReceiverName		- name of the receiver company, your partner/customer name
* IssuerAddress		- self explanatory
* ReceiverAddress	- self explanatory
* IssuerTaxId		- your company tax ID can mean different things in different countries, in eu its VIES no.
* ReceiverTaxId		- your customer's tax ID (if applicable) can mean different things in different countries, in eu its VIES no.
* PaymentType		- how your customer is going to pay you. I usually put "transfer" there.
* DueDate		- expected deadline of money transfer from your customer. In Poland, standard is to put 2 weeks there.
* AccountNo		- self explanatory
* IssuerBankName	- self explanatory
* IssuerBic		- self explanatory

List of names for row of the invoice:

 * Title - product/service name
 * Qty - quantity, number of items
 * Unit - unit based on which quantity is calculated, i.e. 20 kg
 * Price - self explanatory
 * Amount - quantity * price
 * VR - VAT Rate -> integer that is meant to be used as percentage. If you fill in `20` there, that means you have 20% vat rate for that type of goods
 * VA - VAT Amount -> `Price` * `VR`
 * Total - `VA` + `Amount`

Afterwards totals are calculated:
 * ASum - Amount Sum -> sum of all amount values
 * TaxSum - sum of all VAT tax values
 * TotSum - sum of all `Total` values

 You might noticed there are different numbers of spacebars in each field. Also some fields have long names while others have aliases like `VA`. This is because, to fit small numbers into the template that would fit into A4 format, I had to shrink things that I knew will be small such ass Quantity or Vat Rate (most likely 2-digits max). BUT if you want you can give them more space in your config. Example:

> The potential issue you might notice is when there's not enough space for the field to put all the data. This can be dangerous so if you have really big transactions you can give it 8 or 9 digits.

```
[VA   ] // default format if vat amount, will cover 7 digits as square braces are treated as part of the space for the alias.
[Qty] // default format, will cover 5 digits in the quantity column
...
[Qty  ] // modified format, will cover 7 digits in the quantity column! Now you can sell 30 000 pieces of your book!
```

## Formatting

I wanted to keep the format of the document as each shrinked spacebar would destroy the document format (borders wouldn't be aligned). **Invoice-maker treats every excess space as if it meant to be there**. If the value does not fill up whole space it will add additional missing spacebars during replacement. That means it won't destroy document format (each border will be where it meant to be).

Example:

Your vat rate is 7, this is what `invoice-maker` will print instead `[VR]`:

`7  `

This is because sum of text for that alias is 2 for `VR` and two for square braces, ending up one digit and three spaces.

## Print

To print the invoice you have to pick it from the list on the invoice list view and click 'p'. No message will show up as of now but It's on my list!

It will be stored in the folder you specified in the `Config` section. It expects absolute directory so keep that in mind!

## TODO

- [ ] Create proper way to deploy to AUR and other distros. .deb and arch build will come first.
- [ ] Error handling that stops application end returns log in the console.
- [x] Error modal with info what went wrong.
- [ ] Show confirmation of print operation.
- [ ] Figure out how to store $HOME or `~` in the invoice directory. 
- [ ] Write documentation of all the key bindings shown as modal under `H` or `F1`
- [ ] Write man docs or move that README to man format.
