package com.blackmidori.familyexpenses.api;

import org.apache.log4j.BasicConfigurator;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class FamilyExpensesApplication {

	public static void main(String[] args) {
		BasicConfigurator.configure();
		SpringApplication.run(FamilyExpensesApplication.class, args);
	}

}
