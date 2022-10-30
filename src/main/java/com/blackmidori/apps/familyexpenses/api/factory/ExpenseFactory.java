package com.blackmidori.apps.familyexpenses.api.factory;

import com.blackmidori.apps.familyexpenses.api.dto.ExpenseDto;
import com.blackmidori.apps.familyexpenses.api.dto.PayerDto;
import com.blackmidori.apps.familyexpenses.api.model.Expense;
import com.blackmidori.apps.familyexpenses.api.model.Payer;
import com.blackmidori.apps.familyexpenses.api.model.Workspace;
import org.springframework.beans.BeanUtils;

public class ExpenseFactory {

    public Expense createFromDto(ExpenseDto expenseDto, Workspace workspace) {
        Expense expense = new Expense();
        BeanUtils.copyProperties(expenseDto,expense);
        expense.setWorkspace(workspace);
        return expense;
    }
}
