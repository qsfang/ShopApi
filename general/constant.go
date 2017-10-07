/*
 * MIT License
 *
 * Copyright (c) 2017 SmartestEE Inc.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

/*
 * Revision History:
 *     Initial: 2017/07/19        Yusan Kurban
 */

package general

const (
	// General
	// Login session
	SessionUserID   = "userid"
	DuplicateEntry  = "Duplicate"
	InvalidPassword = "match"

	//User
	// User Status
	UserActive   = 0x0
	UserInactive = 0x1

	// sex
	Sex   = 0x0
	Man   = 0x1
	Woman = 0x2

	// Address
	// Address  Default
	AddressNotDefault = 0x0
	AddressDefault    = 0x1

	// Carts
	//CartPro status
	ProInCart    = 0x0
	ProNotInCart = 0x1

	// Category
	//Category Status
	CategoryNotUse = 0x0
	CategoryOnUse  = 0x1

	// Orders
	// Order Status
	OrderUnfinished = 0x0
	OrderFinished   = 0x1
	OrderCanceled   = 0x2
	OrderGetAll     = 0x3

	//Order Payment
	OrderPaid   = 0x1
	OrderUnpaid = 0x2

	// Products
	//Products Status
	ProductOnSale = 0x0
	ProductUnSale = 0x1

	// Product Image Class
	ProductAvatar      = 0x0
	ProductImage       = 0x1
	ProductDetailImage = 0x2

	//Pay
	//Pay Way
	PayOnline  = 0x0
	PayArrive  = 0x1
	PayCompany = 0x2
)
