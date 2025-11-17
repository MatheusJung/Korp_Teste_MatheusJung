using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace InventoryService.Domain.Enums
{
    public enum MovementType
    {
        Entrada,  // Adição de estoque
        Saida,    // Remoção de estoque
        Ajuste    // Ajuste manual de estoque
    }
}
