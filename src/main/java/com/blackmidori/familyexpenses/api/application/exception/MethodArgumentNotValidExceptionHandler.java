package com.blackmidori.familyexpenses.api.application.exception;

import org.springframework.core.Ordered;
import org.springframework.core.annotation.Order;
import org.springframework.validation.BindingResult;
import org.springframework.validation.FieldError;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.ResponseStatus;

import java.time.OffsetDateTime;
import java.time.ZonedDateTime;
import java.util.ArrayList;
import java.util.List;

import static org.springframework.http.HttpStatus.BAD_REQUEST;

@Order(Ordered.HIGHEST_PRECEDENCE)
@ControllerAdvice
public class MethodArgumentNotValidExceptionHandler {

    @ResponseStatus(BAD_REQUEST)
    @ResponseBody
    @ExceptionHandler(MethodArgumentNotValidException.class)
    public Error methodArgumentNotValidException(MethodArgumentNotValidException ex) {
        BindingResult result = ex.getBindingResult();
        List<org.springframework.validation.FieldError> fieldErrors = result.getFieldErrors();
        return processFieldErrors(fieldErrors);
    }

    private Error processFieldErrors(List<org.springframework.validation.FieldError> fieldErrors) {
        Error error = new Error(BAD_REQUEST.value(),  BAD_REQUEST.getReasonPhrase());
        for (org.springframework.validation.FieldError fieldError: fieldErrors) {
            error.addFieldError(fieldError);
        }
        return error;
    }

    static class Error {
        private final OffsetDateTime timestamp;
        private final int status;
        private final String error;
        private final List<FieldError> fieldErrors = new ArrayList<>();

        Error( int status, String error) {
            this.timestamp = OffsetDateTime.now();
            this.status = status;
            this.error = error;
        }
        public OffsetDateTime getTimestamp() {
            return timestamp;
        }

        public int getStatus() {
            return status;
        }

        public String getError() {
            return error;
        }

        public void addFieldError(FieldError fieldError) {
            fieldErrors.add(fieldError);
        }

        public List<FieldError> getFieldErrors() {
            return fieldErrors;
        }
    }
}