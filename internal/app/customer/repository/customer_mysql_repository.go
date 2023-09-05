package repository

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/vamika-digital/wms-api-server/internal/app/customer/domain"
	"github.com/vamika-digital/wms-api-server/internal/utility/customerrors"
	"github.com/vamika-digital/wms-api-server/pkg/database"
)

type MySqlCustomerRepository struct {
	conn database.Connection
}

func NewCustomerRepository(conn database.Connection) CustomerRepository {
	return &MySqlCustomerRepository{conn: conn}
}

func (r *MySqlCustomerRepository) Create(customer *domain.Customer) error {
	query := "INSERT INTO customers (code, name, contact_person, billing_address_address1, billing_address_address2, billing_address_state, billing_address_country, billing_address_pincode, shipping_address_address1, shipping_address_address2, shipping_address_state, shipping_address_country, shipping_address_pincode, status, last_updated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := r.conn.GetDB().Exec(query, customer.Code, customer.Name, customer.ContactPerson, customer.BillingAddress.Address1, customer.BillingAddress.Address2, customer.BillingAddress.State, customer.BillingAddress.Country, customer.BillingAddress.Pincode, customer.ShippingAddress.Address1, customer.ShippingAddress.Address2, customer.ShippingAddress.State, customer.ShippingAddress.Country, customer.ShippingAddress.Pincode, customer.Status, customer.LastUpdatedBy)
	return err
}

func (r *MySqlCustomerRepository) Update(customer *domain.Customer) error {
	query := "UPDATE customers SET code=?, name=?, contact_person=?, billing_address_address1=?, billing_address_address2=?, billing_address_state=?, billing_address_country=?, billing_address_pincode=?, shipping_address_address1=?, shipping_address_address2=?, shipping_address_state=?, shipping_address_country=?, shipping_address_pincode=?, status=?, last_updated_by=? WHERE id=?"
	_, err := r.conn.GetDB().Exec(query, customer.Code, customer.Name, customer.ContactPerson, customer.BillingAddress.Address1, customer.BillingAddress.Address2, customer.BillingAddress.State, customer.BillingAddress.Country, customer.BillingAddress.Pincode, customer.ShippingAddress.Address1, customer.ShippingAddress.Address2, customer.ShippingAddress.State, customer.ShippingAddress.Country, customer.ShippingAddress.Pincode, customer.Status, customer.LastUpdatedBy, customer.ID)
	return err
}

func (r *MySqlCustomerRepository) Delete(customerID int64) error {
	query := "DELETE FROM customers WHERE id=?"
	_, err := r.conn.GetDB().Exec(query, customerID)
	return err
}

func (r *MySqlCustomerRepository) GetById(customerID int64) (*domain.Customer, error) {
	query := "SELECT id, code, name, contact_person, billing_address_address1, billing_address_address2, billing_address_state, billing_address_country, billing_address_pincode, shipping_address_address1, shipping_address_address2, shipping_address_state, shipping_address_country, shipping_address_pincode, status, created_at, updated_at, last_updated_by FROM customers WHERE id = ?"
	row := r.conn.GetDB().QueryRow(query, customerID)
	customer := domain.NewCustomerWithDefaults()
	err := row.Scan(&customer.ID, &customer.Code, &customer.Name, &customer.ContactPerson, &customer.BillingAddress.Address1, &customer.BillingAddress.Address2, &customer.BillingAddress.State, &customer.BillingAddress.Country, &customer.BillingAddress.Pincode, &customer.ShippingAddress.Address1, &customer.ShippingAddress.Address2, &customer.ShippingAddress.State, &customer.ShippingAddress.Country, &customer.ShippingAddress.Pincode, &customer.Status, &customer.CreatedAt, &customer.UpdatedAt, &customer.LastUpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerrors.ErrResourceNotFound
		}
		return nil, err
	}
	return customer, nil
}

func (r *MySqlCustomerRepository) GetTotalCount(filter CustomerFilterOptions) (int, error) {
	query, args := r.buildFilterQuery("SELECT COUNT(*) FROM customers", filter)

	var count int
	if err := r.conn.GetDB().QueryRow(query, args...).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *MySqlCustomerRepository) GetAll(page int, pageSize int, sort string, filter CustomerFilterOptions) ([]*domain.Customer, error) {
	query, args := r.buildFilterQuery("SELECT id, code, name, contact_person, billing_address_address1, billing_address_address2, billing_address_state, billing_address_country, billing_address_pincode, shipping_address_address1, shipping_address_address2, shipping_address_state, shipping_address_country, shipping_address_pincode, status, created_at, updated_at, last_updated_by FROM customers", filter)
	var allowedSortOrders = map[string]bool{
		"code ASC":            true,
		"code DESC":           true,
		"name ASC":            true,
		"name DESC":           true,
		"contact_person ASC":  true,
		"contact_person DESC": true,
		"status ASC":          true,
		"status DESC":         true,
	}

	if sort != "" {
		if _, ok := allowedSortOrders[sort]; ok {
			query += " ORDER BY " + sort
		} else {
			return nil, errors.New("invalid sort order")
		}
	}

	query += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, page)

	rows, err := r.conn.GetDB().Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []*domain.Customer
	for rows.Next() {
		customer := domain.NewCustomerWithDefaults()
		if err := rows.Scan(&customer.ID, &customer.Code, &customer.Name, &customer.ContactPerson, &customer.BillingAddress.Address1, &customer.BillingAddress.Address2, &customer.BillingAddress.State, &customer.BillingAddress.Country, &customer.BillingAddress.Pincode, &customer.ShippingAddress.Address1, &customer.ShippingAddress.Address2, &customer.ShippingAddress.State, &customer.ShippingAddress.Country, &customer.ShippingAddress.Pincode, &customer.Status, &customer.CreatedAt, &customer.UpdatedAt, &customer.LastUpdatedBy); err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}
	return customers, nil
}

func (r *MySqlCustomerRepository) buildFilterQuery(baseQuery string, filter CustomerFilterOptions) (string, []interface{}) {
	var filters []string
	var args []interface{}

	if filter.Name != "" {
		filters = append(filters, "name LIKE ?")
		args = append(args, "%"+filter.Name+"%")
	}
	if filter.Code != "" {
		filters = append(filters, "code LIKE ?")
		args = append(args, "%"+filter.Code+"%")
	}
	if filter.Status != "" {
		filters = append(filters, "status = ?")
		args = append(args, filter.Status)
	}

	if len(filters) > 0 {
		return baseQuery + " WHERE " + strings.Join(filters, " AND "), args
	}

	return baseQuery, args
}
