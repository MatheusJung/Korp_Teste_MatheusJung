using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace InventoryService.Application.Exceptions
{
    public class OutOfStockException : Exception
    {
        public OutOfStockException(string productCode, int requested, int available)
            : base($"Not enough stock for product '{productCode}'. Requested: {requested}, Available: {available}.")
        {
        }
    }
}
