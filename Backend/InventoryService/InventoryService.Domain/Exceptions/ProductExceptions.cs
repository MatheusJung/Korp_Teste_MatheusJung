namespace InventoryService.Domain.Exceptions
{
    public class ProductException : Exception
    {
        public ProductException(string message) : base(message) { }
    }

    public static class ProductExceptions
    {
        public static ProductException DuplicateCode(string code) =>
            new ProductException($"Product with code '{code}' already exists.");
        public static ProductException InsufficientStock(string code, int requested, int available) =>
            new ProductException($"Not enough stock for product '{code}'. Requested: {requested}, Available: {available}.");
        public static ProductException AlreadyDeactivated(string code) =>
            new ProductException($"Product '{code}' is already deactivated.");
        public static ProductException AlreadyActivated(string code) =>
            new ProductException($"Product '{code}' is already activated.");
        public static ProductException NotFound(string code) =>
            new ProductException($"Product '{code}' not found.");
        public static ProductException InvalidCode(string code) =>
            new ProductException($"Product code '{code}' is invalid. Only uppercase letters, numbers, . and - are allowed, no spaces.");
        public static ProductException InvalidName(string name) =>
            new ProductException($"Product name '{name}' is invalid. Only letters, numbers and spaces are allowed.");
    }
}