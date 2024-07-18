
# E-Commerce Service

This service provides basic e-commerce functionalities including viewing products, managing shopping carts, checking out, and user authentication.

## Features

- **View Products by Category**: Customers can view a list of products filtered by category.
- **Add Products to Shopping Cart**: Customers can add products to their shopping cart.
- **View Shopping Cart**: Customers can see a list of products that have been added to their shopping cart.
- **Delete Products from Shopping Cart**: Customers can delete products from their shopping cart.
- **Checkout and Payment**: Customers can checkout and make payment transactions.
- **User Authentication**: Customers can register and login.

## API Documentation

The API documentation is available at the following route:
```
/swagger
```

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Installation

1. **Clone the repository**

   ```bash
   git clone https://github.com/your-repository.git
   cd your-repository
   ```

2. **Copy the example environment file**

   ```bash
   cp .env.example .env
   ```

3. **Run the application**

   ```bash
   make dev
   ```

## Environment Variables

Ensure you have the following environment variables set in your `.env` file:

```
# Database Configuration
DB_HOST=your_db_host
DB_PORT=your_db_port
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name

# JWT Secret
JWT_SECRET=your_jwt_secret

# Other configurations
OTHER_CONFIG=your_other_config
```

### Example `.env` file

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=ecommerce

# JWT Secret
JWT_SECRET=your_jwt_secret

# Other configurations
OTHER_CONFIG=value
```

## Docker

The Docker image for this service is available on Docker Hub: `azka1415/synap-be`

To pull and run the Docker image:

```bash
docker pull azka1415/synap-be
docker run -d -p 8080:8080 --env-file .env azka1415/synap-be
```

## Usage

Once the service is up and running, you can access the API documentation at `http://localhost:8080/swagger`.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any changes.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.
