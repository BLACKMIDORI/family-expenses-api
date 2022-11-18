package com.blackmidori.familyexpenses.api.application;

import com.blackmidori.familyexpenses.api.application.converter.OffsetDateTimeReadConverter;
import com.blackmidori.familyexpenses.api.application.converter.OffsetDateTimeWriteConverter;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.mongodb.core.convert.MongoCustomConversions;

import java.util.Arrays;

@Configuration
public class SpringConfiguration {

  @Bean
  public MongoCustomConversions mongoCustomConversions() {

    return new MongoCustomConversions(
        Arrays.asList(
                new OffsetDateTimeWriteConverter(),
                new OffsetDateTimeReadConverter()
        ));
  }

}