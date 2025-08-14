IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'Products' AND type = 'U')
BEGIN
    CREATE TABLE Products (
        ID INT IDENTITY(1,1) PRIMARY KEY,
        Name NVARCHAR(255) NOT NULL,
        Description NVARCHAR(MAX) NOT NULL,
        Image VARCHAR(255) NOT NULL,
        Price DECIMAL(10, 2) NOT NULL CHECK (Price >= 0),
        Quantity INT NOT NULL CHECK (Quantity >= 0) ,
        CreatedAt DATETIME2 NOT NULL DEFAULT GETDATE()
    );
END
