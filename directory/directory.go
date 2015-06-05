// The directory package includes the interface for the Otto Appfile
// directory service that stores data related to Appfiles.
//
// The directory service is used to maintain forward and backward
// indexes of applications to projects, projects to infras, etc. This
// allows Otto to gain a global view of an Appfile when it is being
// used.
//
// The existence of such a service allows Appfiles to not have to contain
// global information of a potentially complex project.
package directory
