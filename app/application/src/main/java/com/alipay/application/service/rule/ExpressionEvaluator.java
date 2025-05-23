/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package com.alipay.application.service.rule;

import java.util.Map;

public class ExpressionEvaluator {
  
    public static boolean evalExpression(String condition, Map<Integer, Boolean> resultsMap) {
        condition = condition.replaceAll("\\s+", "");
        return evaluate(condition, resultsMap);
    }

    private static boolean evaluate(String condition, Map<Integer, Boolean> resultsMap) {
        if (condition == null || condition.isEmpty()) {
            throw new IllegalArgumentException("Condition is empty");
        }
        return evaluateExpression(condition, resultsMap, 0).result;
    }

    private static EvalResult evaluateExpression(String condition, Map<Integer, Boolean> resultsMap, int index) {
        boolean currentResult = true;
        boolean currentLogicalAnd = true;
        int i = index;

        while (i < condition.length()) {
            char currentChar = condition.charAt(i);
            if (Character.isDigit(currentChar)) {
                int num = currentChar - '0';
                boolean value = resultsMap.getOrDefault(num, false);
                if (currentLogicalAnd) {
                    currentResult = currentResult && value;
                } else {
                    currentResult = currentResult || value;
                }
                i++;
                
            } else if (currentChar == '&') {
                i += 2; // Skip "&&"
                currentLogicalAnd = true;
            } else if (currentChar == '|') {
                i += 2; // Skip "||"
                currentLogicalAnd = false;
            } else if (currentChar == '(') {
                EvalResult result = evaluateExpression(condition, resultsMap, i + 1);
                if (currentLogicalAnd) {
                    currentResult = currentResult && result.result;
                } else {
                    currentResult = currentResult || result.result;
                }
                i = result.index + 1;
            } else if (currentChar == ')') {
                return new EvalResult(currentResult, i);
            } else {
                throw new IllegalArgumentException("Invalid character in condition: " + currentChar);
            }
        }

        return new EvalResult(currentResult, i);
    }

    private static class EvalResult {
        boolean result;
        int index;

        EvalResult(boolean result, int index) {
            this.result = result;
            this.index = index;
        }
    }

//    public static void main(String[] args) {
//        String condition = "1 && 2";
//        Map<Integer, Boolean> resultsMap = Map.of(1, true, 2, false, 3, true);
//
//        boolean result = evalExpression(condition, resultsMap);
//        System.out.println("Result: " + result);
//    }
}
