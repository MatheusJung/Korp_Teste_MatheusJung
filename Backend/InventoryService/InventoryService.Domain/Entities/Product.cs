using InventoryService.Domain.Exceptions;
using System.Text.RegularExpressions;

namespace InventoryService.Domain.Entities;

public sealed class Product
{
    public int Id { get; private set; }
    public string Code { get; private set; } = string.Empty; // novo: código externo
    public string Name { get; private set; } = string.Empty;
    public int Quantity { get; private set; }
    public bool IsActive { get; private set; } = true;
    public DateTime CreatedAt { get; private set; } = DateTime.UtcNow;
    public DateTime? DeactivatedAt { get; private set; }

    private Product() { } // EF Core

    public Product(string code, string name, int quantity)
    {
        if (string.IsNullOrWhiteSpace(code))
            throw new ArgumentException("Code is required.", nameof(code));
        if (!IsValidCode(code))
            throw ProductExceptions.InvalidCode(code);
        if (string.IsNullOrWhiteSpace(name))
            throw new ArgumentException("Name is required.", nameof(name));
        if (!IsValidName(name))
            throw ProductExceptions.InvalidName(name);
        if (quantity < 0)
            throw new ArgumentOutOfRangeException(nameof(quantity), "Quantity cannot be negative.");

        Code = code;
        Name = name;
        Quantity = quantity;
        CreatedAt = DateTime.UtcNow;
        IsActive = true;
    }

    private bool IsValidCode(string code) =>
        Regex.IsMatch(code, @"^[A-Z0-9.-]+$");

    private bool IsValidName(string name) =>
        Regex.IsMatch(name, @"^[A-Za-z0-9 ]+$");

    public void DecreaseStock(int amount)
    {
        if (amount <= 0)
            throw new ArgumentOutOfRangeException(nameof(amount), "Amount must be greater than zero.");
        if (Quantity < amount)
            throw ProductExceptions.InsufficientStock(Code, amount, Quantity);
        Quantity -= amount;
    }

    public void IncreaseStock(int amount)
    {
        if (amount <= 0) 
            throw new ArgumentOutOfRangeException(nameof(amount), "Amount must be greater than zero.");
        
        Quantity += amount;
    }

    public void Deactivate()
    {
        if (!IsActive)
            throw ProductExceptions.AlreadyDeactivated(Code);
        IsActive = false;
        DeactivatedAt = DateTime.UtcNow;
    }

    public void Activate()
    {
        if (IsActive)
            throw ProductExceptions.AlreadyActivated(Code);
        IsActive = true;
        DeactivatedAt = null;
    }
}
