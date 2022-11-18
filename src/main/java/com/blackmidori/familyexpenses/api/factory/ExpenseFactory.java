package com.blackmidori.familyexpenses.api.factory;

import com.blackmidori.familyexpenses.api.dto.ExpenseDto;
import com.blackmidori.familyexpenses.api.model.Expense;
import com.blackmidori.familyexpenses.api.model.Workspace;
import org.springframework.beans.BeanUtils;

public class ExpenseFactory {

    public Expense createFromDto(ExpenseDto expenseDto, Workspace workspace) {
        Expense expense = new Expense();
        BeanUtils.copyProperties(expenseDto,expense);
        expense.setWorkspace(workspace);
        return expense;
    }
}
