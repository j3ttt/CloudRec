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
package com.alipay.application.service.risk;

import com.alipay.application.share.vo.ApiResponse;
import com.alipay.application.share.vo.ListVO;
import com.alipay.application.share.vo.rule.RuleScanResultVO;
import com.alipay.dao.dto.RuleScanResultDTO;
import com.alipay.dao.dto.RuleStatisticsDTO;
import jakarta.servlet.http.HttpServletResponse;

import java.io.IOException;
import java.util.List;

/*
 *@title RiskService
 *@description
 *@author jietian
 *@version 1.0
 *@create 2024/7/16 16:33
 */
public interface RiskService {

    ApiResponse<ListVO<RuleScanResultVO>> queryRiskList(RuleScanResultDTO dto);


    ApiResponse<RuleScanResultVO> queryRiskDetail(Long riskId);


    ApiResponse<String> ignoreRisk(Long riskId, String ignoreReason, String ignoreReasonType);


    ApiResponse<String> cancelIgnoreRisk(RuleScanResultDTO dto);

    void exportRiskList(HttpServletResponse response, RuleScanResultDTO dto) throws IOException;

    List<RuleStatisticsDTO> listRuleStatistics(RuleScanResultDTO ruleScanResultDTO);
}
