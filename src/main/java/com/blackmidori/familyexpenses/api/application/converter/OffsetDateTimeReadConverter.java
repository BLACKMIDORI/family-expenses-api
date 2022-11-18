package com.blackmidori.familyexpenses.api.application.converter;

import org.apache.log4j.helpers.ISO8601DateFormat;
import org.springframework.core.convert.converter.Converter;

import java.time.OffsetDateTime;
import java.time.format.DateTimeFormatter;

public class OffsetDateTimeReadConverter implements Converter<String, OffsetDateTime> {
    @Override
    public OffsetDateTime convert(String offsetDateTime) {

        return OffsetDateTime.parse(offsetDateTime);
    }

}