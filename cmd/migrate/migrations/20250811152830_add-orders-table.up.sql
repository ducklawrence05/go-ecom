IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'Orders' AND type = 'U')
BEGIN
    CREATE TABLE Orders (
        ID INT IDENTITY(1,1) PRIMARY KEY,
        UserID INT NOT NULL,
        Total DECIMAL(10, 2) NOT NULL CHECK (Total >= 0),
        Status NVARCHAR(50) NOT NULL CHECK (Status IN ('pending', 'completed', 'cancelled')),
        Address NVARCHAR(255) NOT NULL,
        CreatedAt DATETIME2 NOT NULL DEFAULT GETDATE(),
    
        CONSTRAINT FK_Orders_Users FOREIGN KEY (UserID) REFERENCES Users(ID) ON DELETE CASCADE    
    );
END