namespace BillingService.Application.Exceptions
{
    public class OperationAlreadyAppliedException : Exception
    {
        public string OperationKey { get; }

        public OperationAlreadyAppliedException(string operationKey)
            : base($"Operation with key '{operationKey}' has already been applied.")
        {
            OperationKey = operationKey;
        }
    }
}
