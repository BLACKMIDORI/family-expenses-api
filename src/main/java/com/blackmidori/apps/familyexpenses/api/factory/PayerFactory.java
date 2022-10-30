package com.blackmidori.apps.familyexpenses.api.factory;

import com.blackmidori.apps.familyexpenses.api.dto.PayerDto;
import com.blackmidori.apps.familyexpenses.api.model.Payer;
import com.blackmidori.apps.familyexpenses.api.model.Workspace;
import org.springframework.beans.BeanUtils;

import java.time.OffsetDateTime;

public class PayerFactory {

    public Payer createFromDto(PayerDto payerDto, Workspace workspace) {
        Payer payer = new Payer();
        BeanUtils.copyProperties(payerDto,payer);
        payer.setWorkspace(workspace);
        return payer;
    }
}
