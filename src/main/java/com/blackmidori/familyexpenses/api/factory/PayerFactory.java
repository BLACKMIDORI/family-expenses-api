package com.blackmidori.familyexpenses.api.factory;

import com.blackmidori.familyexpenses.api.dto.PayerDto;
import com.blackmidori.familyexpenses.api.model.Payer;
import com.blackmidori.familyexpenses.api.model.Workspace;
import org.springframework.beans.BeanUtils;

public class PayerFactory {

    public Payer createFromDto(PayerDto payerDto, Workspace workspace) {
        Payer payer = new Payer();
        BeanUtils.copyProperties(payerDto,payer);
        payer.setWorkspace(workspace);
        return payer;
    }
}
