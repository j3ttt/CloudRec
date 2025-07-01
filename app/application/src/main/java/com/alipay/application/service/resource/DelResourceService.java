package com.alipay.application.service.resource;


/*
 *@title DelResourceService
 *@description
 *@author suitianshuang
 *@version 1.0
 *@create 2025/6/26 22:17
 */
public interface DelResourceService {

    /**
     * 预删除资源，将资源的逻辑删除次数 + 1
     *
     * @param cloudAccountId
     * @return
     */
    int preDeleteByCloudAccountId(String cloudAccountId);

}
