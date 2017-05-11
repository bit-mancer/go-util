package logging

// SourceKey is the map key used for logging sources (see NewLog).
const SourceKey = "_source"

// DomainKey is the map key used for logging domains (see NewDomainLogger).
const DomainKey = "_domain"

// HostKey is the map key used for logging the machine hostname (via os.Hostname).
const HostKey = "_host"
