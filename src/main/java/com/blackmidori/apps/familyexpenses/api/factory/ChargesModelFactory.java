package com.blackmidori.apps.familyexpenses.api.factory;

import com.blackmidori.apps.familyexpenses.api.application.exception.EntityNotFound;
import com.blackmidori.apps.familyexpenses.api.dto.ChargesModelDto;
import com.blackmidori.apps.familyexpenses.api.model.*;
import com.blackmidori.apps.familyexpenses.api.repository.ExpenseRepository;
import com.blackmidori.apps.familyexpenses.api.repository.PayerRepository;
import com.blackmidori.apps.familyexpenses.api.repository.WorkspaceRepository;
import org.springframework.beans.BeanUtils;
import org.springframework.stereotype.Component;

import java.util.Optional;

@Component
public class ChargesModelFactory {

    private final WorkspaceRepository workspaceRepository;
    private final ExpenseRepository expenseRepository;
    private final PayerRepository payerRepository;

    public ChargesModelFactory(WorkspaceRepository workspaceRepository, ExpenseRepository expenseRepository, PayerRepository payerRepository) {
        this.workspaceRepository = workspaceRepository;
        this.expenseRepository = expenseRepository;
        this.payerRepository = payerRepository;
    }

    public ChargesModel createFromDto(ChargesModelDto chargesModelDto) throws EntityNotFound {
        ChargesModel chargesModel = new ChargesModel();
        // Copy flat properties
        BeanUtils.copyProperties(chargesModelDto,chargesModel);
        // Load other instances
        Optional<Workspace> workspaceOptional = workspaceRepository.findById(chargesModelDto.getWorkspaceId());
        if(workspaceOptional.isEmpty()){
            throw new EntityNotFound(Workspace.class, chargesModelDto.getWorkspaceId());
        }
        chargesModel.setWorkspace(workspaceOptional.get());
        // Manual set list
        chargesModel.setChargesAssociations(chargesModelDto.getChargesAssociations().stream()
                .map(chargeAssociationDto -> {
                    var chargeAssociation = new ChargeAssociation();
                    // Copy flat properties
                    BeanUtils.copyProperties(chargeAssociationDto,chargeAssociation);
                    // Load other instances
                    final var expenseOptional = expenseRepository.findById(chargeAssociationDto.getExpenseId());
                    if(expenseOptional.isEmpty()){
                        throw new EntityNotFound(Expense.class,chargeAssociationDto.getExpenseId());
                    }
                    chargeAssociation.setExpense(expenseOptional.get());
                    // Manual set list
                    chargeAssociation.setPaymentWeights(chargeAssociationDto.getPaymentWeights().stream()
                            .map(payerPaymentWeightDto -> {
                                var payerPaymentWeight = new PayerPaymentWeight();
                                BeanUtils.copyProperties(payerPaymentWeightDto,payerPaymentWeight);
                                final var payerOptional = payerRepository.findById(payerPaymentWeightDto.getPayerId());
                                if(payerOptional.isEmpty()){
                                    throw new EntityNotFound(Payer.class,payerPaymentWeightDto.getPayerId());
                                }
                                payerPaymentWeight.setPayer(payerOptional.get());
                                return payerPaymentWeight;
                            }).toList()
                    );
                    return chargeAssociation;
                }).toList()
        );
        return chargesModel;
    }
}
