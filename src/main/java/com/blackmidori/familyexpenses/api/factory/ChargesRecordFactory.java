package com.blackmidori.familyexpenses.api.factory;

import com.blackmidori.familyexpenses.api.application.exception.EntityNotFound;
import com.blackmidori.familyexpenses.api.dto.ChargesRecordDto;
import com.blackmidori.familyexpenses.api.model.*;
import com.blackmidori.familyexpenses.api.repository.ChargesModelRepository;
import com.blackmidori.familyexpenses.api.repository.ExpenseRepository;
import com.blackmidori.familyexpenses.api.repository.PayerRepository;
import com.blackmidori.familyexpenses.api.repository.WorkspaceRepository;
import org.springframework.beans.BeanUtils;
import org.springframework.stereotype.Component;

import java.util.Optional;

@Component
public class ChargesRecordFactory {

    private final WorkspaceRepository workspaceRepository;
    private final ExpenseRepository expenseRepository;
    private final PayerRepository payerRepository;
    private final ChargesModelRepository chargesModelRepository;

    public ChargesRecordFactory(WorkspaceRepository workspaceRepository, ExpenseRepository expenseRepository, PayerRepository payerRepository, ChargesModelRepository chargesModelRepository) {
        this.workspaceRepository = workspaceRepository;
        this.expenseRepository = expenseRepository;
        this.payerRepository = payerRepository;
        this.chargesModelRepository = chargesModelRepository;
    }

    public ChargesRecord createFromDto(ChargesRecordDto chargesRecordDto) throws EntityNotFound {
        ChargesRecord chargesRecord = new ChargesRecord();
        // Copy flat properties
        BeanUtils.copyProperties(chargesRecordDto,chargesRecord);
        // Load other instances
        Optional<ChargesModel> chargesModelOptional = chargesModelRepository.findById(chargesRecordDto.getChargesModelId());
        if(chargesModelOptional.isEmpty()){
            throw new EntityNotFound(ChargesModel.class, chargesRecordDto.getChargesModelId());
        }
        chargesRecord.setChargesModel(chargesModelOptional.get());
        // Manual set list
        chargesRecord.setCharges(chargesRecordDto.getCharges().stream()
                .map(chargeDto -> {
                    var charge = new Charge();
                    // Copy flat properties
                    BeanUtils.copyProperties(chargeDto,charge);
                    // Manual set complex properties
                    var bill = new Bill();
                    BeanUtils.copyProperties(chargeDto.getBill(),bill);
                    charge.setBill(bill);
                    // Load other instances
                    Optional<Expense> expenseOptional = expenseRepository.findById(chargeDto.getBill().getExpenseId());
                    if(expenseOptional.isEmpty()){
                        throw new EntityNotFound(Expense.class, chargeDto.getBill().getExpenseId());
                    }
                    bill.setExpense(expenseOptional.get());
                    // Manual set list
                    charge.setPaymentAmountList(chargeDto.getPaymentAmountList().stream()
                            .map(payerPaymentAmountDto -> {
                                var payerPaymentAmount = new PayerPaymentAmount();
                                BeanUtils.copyProperties(payerPaymentAmountDto,payerPaymentAmount);
                                final var payerOptional = payerRepository.findById(payerPaymentAmountDto.getPayerId());
                                if(payerOptional.isEmpty()){
                                    throw new EntityNotFound(Payer.class,payerPaymentAmountDto.getPayerId());
                                }
                                payerPaymentAmount.setPayer(payerOptional.get());
                                return payerPaymentAmount;
                            }).toList()
                    );
                    return charge;
                }).toList()
        );
        return chargesRecord;
    }
}
