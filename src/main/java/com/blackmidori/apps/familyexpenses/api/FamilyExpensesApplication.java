package com.blackmidori.apps.familyexpenses.api;

import org.apache.log4j.BasicConfigurator;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.MessageSource;
import org.springframework.context.annotation.Bean;
import org.springframework.context.support.ReloadableResourceBundleMessageSource;
import org.springframework.validation.beanvalidation.LocalValidatorFactoryBean;

@SpringBootApplication
public class FamilyExpensesApplication {

	public static void main(String[] args) {
		BasicConfigurator.configure();
		SpringApplication.run(FamilyExpensesApplication.class, args);
	}

}
