package rds_test

const (
	// Please make sure GovCloud and commercial support these since they vary
	postgresPreferredInstanceClasses    = `"db.t3.micro", "db.t3.small", "db.t2.small", "db.t2.medium"`
	mySQLPreferredInstanceClasses       = `"db.t3.micro", "db.t3.small", "db.t2.small", "db.t2.medium"`
	mariaDBPreferredInstanceClasses     = `"db.t3.micro", "db.t3.small", "db.t2.small", "db.t2.medium"`
	oraclePreferredInstanceClasses      = `"db.t3.medium", "db.t2.medium", "db.t3.large", "db.t2.large"` // Oracle requires at least a medium instance as a replica source
	sqlServerPreferredInstanceClasses   = `"db.t2.small", "db.t3.small"`
	sqlServerSEPreferredInstanceClasses = `"db.m5.large", "db.m4.large", "db.r4.large"`
	oracleSE2PreferredInstanceClasses   = `"db.m5.large", "db.m4.large", "db.r4.large"`
)
