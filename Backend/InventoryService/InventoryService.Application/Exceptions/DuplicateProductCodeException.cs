using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace InventoryService.Application.Exceptions
{
    public class DuplicateProductCodeException : Exception
    {
        public DuplicateProductCodeException(string code)
    :       base($"Product with code '{code}' already exists.") { }
    }
}
