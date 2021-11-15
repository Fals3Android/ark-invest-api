package main

import "github.com/aws/aws-sdk-go/service/dynamodb"

type MetaInfo struct {
	url             string
	tableName       string
	data            [][]string
	attributeValues []map[string]*dynamodb.AttributeValue
}

func getTablesAndUrls() map[string]MetaInfo {
	info := make(map[string]MetaInfo)
	info["ARKK"] = MetaInfo{
		url:       "https://ark-funds.com/wp-content/uploads/funds-etf-csv/ARK_INNOVATION_ETF_ARKK_HOLDINGS.csv",
		tableName: "ARK_INNOVATION_ETF_ARKK_HOLDINGS",
	}
	info["ARKW"] = MetaInfo{
		url:       "https://ark-funds.com/wp-content/uploads/funds-etf-csv/ARK_NEXT_GENERATION_INTERNET_ETF_ARKW_HOLDINGS.csv",
		tableName: "ARK_NEXT_GENERATION_INTERNET_ETF_ARKW_HOLDINGS",
	}
	info["ARKG"] = MetaInfo{
		url:       "https://ark-funds.com/wp-content/uploads/funds-etf-csv/ARK_GENOMIC_REVOLUTION_ETF_ARKG_HOLDINGS.csv",
		tableName: "ARK_GENOMIC_REVOLUTION_ETF_ARKG_HOLDINGS",
	}
	info["ARKF"] = MetaInfo{
		url:       "https://ark-funds.com/wp-content/uploads/funds-etf-csv/ARK_FINTECH_INNOVATION_ETF_ARKF_HOLDINGS.csv",
		tableName: "ARK_FINTECH_INNOVATION_ETF_ARKF_HOLDINGS",
	}
	info["ARKQ"] = MetaInfo{
		url:       "https://ark-funds.com/wp-content/uploads/funds-etf-csv/ARK_AUTONOMOUS_TECH._&_ROBOTICS_ETF_ARKQ_HOLDINGS.csv",
		tableName: "ARK_AUTONOMOUS_TECH_AND_ROBOTICS_ETF_ARKQ_HOLDINGS",
	}
	info["ARKX"] = MetaInfo{
		url:       "https://ark-funds.com/wp-content/uploads/funds-etf-csv/ARK_SPACE_EXPLORATION_&_INNOVATION_ETF_ARKX_HOLDINGS.csv",
		tableName: "ARK_SPACE_EXPLORATION_AND_INNOVATION_ETF_ARKX_HOLDINGS",
	}
	info["PRNT"] = MetaInfo{
		url:       "https://ark-funds.com/wp-content/uploads/funds-etf-csv/THE_3D_PRINTING_ETF_PRNT_HOLDINGS.csv",
		tableName: "THE_3D_PRINTING_ETF_PRNT_HOLDINGS",
	}
	info["IZRL"] = MetaInfo{
		url:       "https://ark-funds.com/wp-content/uploads/funds-etf-csv/ARK_ISRAEL_INNOVATIVE_TECHNOLOGY_ETF_IZRL_HOLDINGS.csv",
		tableName: "ARK_ISRAEL_INNOVATIVE_TECHNOLOGY_ETF_IZRL_HOLDINGS",
	}
	return info
}
