package com.blackmidori.familyexpenses.api;

import com.blackmidori.familyexpenses.api.service.PayerService;
import com.blackmidori.familyexpenses.api.service.WorkspaceService;
import org.junit.After;
import org.junit.Before;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.ActiveProfiles;
import org.springframework.test.context.TestExecutionListeners;
import org.springframework.test.context.support.DependencyInjectionTestExecutionListener;

@SpringBootTest
@TestExecutionListeners({ DependencyInjectionTestExecutionListener.class })
@ActiveProfiles("test")
public class PayerJavaApiUnitTests {

    @Autowired
    private PayerService payerService;

    @Autowired
    private WorkspaceService workspaceService;

    @Before
    public void setUp() {
        // TODO
    }

    @Test
    public void shouldReturnNotNullPayerService() {
        Assertions.assertNotNull(payerService);
    }

    @Test
    public void shouldReturnNotNullWorkspaceService() {
        Assertions.assertNotNull(workspaceService);
    }

    @Test
    public void shouldReturnTODOTODO() {
        // TODO
    }

    @After
    public void tearDown() {
        // TODO
    }
}