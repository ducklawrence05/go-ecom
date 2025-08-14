IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'OrderItems' AND type = 'U')
BEGIN
    CREATE TABLE OrderItems (
        ID INT IDENTITY(1,1) PRIMARY KEY,
        OrderID INT NOT NULL,
        ProductID INT NOT NULL,
        Price DECIMAL(10, 2) NOT NULL CHECK (Price >= 0),
        Quantity INT NOT NULL CHECK (Quantity >= 0) ,
        CreatedAt DATETIME2 NOT NULL DEFAULT GETDATE(),
    
        UNIQUE (OrderID, ProductID),
        CONSTRAINT FK_OrderItems_Orders FOREIGN KEY (OrderID) REFERENCES Orders(ID) ON DELETE CASCADE,
        CONSTRAINT FK_OrderItems_Products FOREIGN KEY (ProductID) REFERENCES Products(ID) ON DELETE CASCADE
    );
END