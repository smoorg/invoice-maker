% IM(1) im 0.0.1

# NAME

im - invoice-maker, a terminal user interface program to create and manage invoices.

# SYNOPSIS

**im** [*-C*] 

# DESCRIPTION

**im** is a go program that's meant to manage company data and past invoices as well as save them in specified directory.

# OPTIONS

**-c** _dir_, **-C** _dir_, **\-\-config-dir=_dir_**
: Specifies config file directory. By default we use `$XDG_CONFIG_HOME/invoice-maker`. Files like `config.yml`, `template.md` and `template-row.md` are required for application to run smoothly.

**-i** _dir_, **-I** _dir_, **\-\-invoice-dir=_dir_**
: Specifies invoices directory. By default we use `$XDG_DATA_HOME/invoice-maker`. Folder structure inside it has been arbitrary set to `year/month/your_invoice_name.pdf` but that might be a subject for change in future.

# FILES

**im** uses following files and they are expected to land in `$XDG_CONFIG_HOME` (unless -c option says otherwise). If any of these files is not there, it will be generated from defaults stored most likely at `/etc/invoice-maker`.

**config.yml**
: Stores all the data about past receipent and the issuer to make generation process just a dropdown pick. It also stores all the generated invoices for historical purposes in case you loose your files. Made that way to be able to archive the memory and move between machines easily. More about config structure in section five **im(5)**

**template.md**
: Defines main invoice template. You can scrap out things you don't need (i.e. your country has no vat rate for your services and it serves no purpose) or adjust it for your needs. The most important part are square braces which define variable names and space left to print them. Overall on very early stages it was assumed that the template should not be modified in terms of number of characters per line, therefore if you keep too little space for variables, they will be cut off. More about templates in section five **im(5)**.

**template-row.md**
: Same as `template.md`, adjustable definition of a single row of an invoice. More about templates in section five **im(5)**.

# REPORTING BUGS

Report bugs at:

- https://github.com/smoorg/invoice-maker/issues 
- mail it over to mateusz@mateuszreszka.xyz

# COPYRIGHT

Copyright (C) 2023 Mateusz Reszka

**im** is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

**im** is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with **im**.  If not, see <http://www.gnu.org/licenses/>.

# AUTHOR

Mateusz Reszka <mateusz@mateuszreszka.xyz>

For more information regarding the program, see invoice-maker's README at 
https://github.com/smoorg/invoice-maker/blob/master/README.md

