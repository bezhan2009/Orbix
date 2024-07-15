// opm.h
#ifndef OPM_H
#define OPM_H

void show_help();
void install_package(const char *package_name);
void remove_package(const char *package_name);
void update_package_manager();
void list_installed_packages();

#endif // OPM_H
