%ifarch aarch64
%global _arch aarch64
%global BuildArchitectures aarch64
%endif

%ifarch x86_64
%global _arch x86_64
%global BuildArchitectures x86_64
%endif

%define debug_package   %{nil}
%define _build_id_links none
%define _name pgtools
%define _prefix /opt
%define _version 1.70.00
%define _rel 0
#%define _arch x86_64
%define _binaryname pgtools

Name:       pgtools
Version:    %{_version}
Release:    %{_rel}
Summary:    PGSQL multi-usage tool

Group:      PGSQL utils
License:    GPL2.0
URL:        https://git.famillegratton.net:3000/devops/pgtools.git

Source0:    %{name}-%{_version}.tar.gz
#BuildArchitectures: x86_64
BuildRequires: gcc
%description
PGSQL multi-usage tool

%prep
%autosetup

%build
cd %{_sourcedir}/%{_name}-%{_version}/src
PATH=$PATH:/opt/go/bin CGO_ENABLED=0 go build -o %{_sourcedir}/%{_binaryname} .
strip %{_sourcedir}/%{_binaryname}

%clean
rm -rf $RPM_BUILD_ROOT

%pre
exit 0

%install
install -Dpm 0755 %{_sourcedir}/%{_binaryname} %{buildroot}%{_bindir}/%{_binaryname}

%post

%preun

%postun

%files
%defattr(-,root,root,-)
%{_bindir}/%{_binaryname}


%changelog
* Fri Sep 12 2025 Binary package builder <builder@famillegratton.net> 1.50.00-0
- Completed the db create/drop subcomments (jean-francois@famillegratton.net)

* Thu Sep 11 2025 Binary package builder <builder@famillegratton.net> 1.40.00-0
- Fixed flag display in srv -h, doc update (jean-francois@famillegratton.net)
- version bump and fix in env rm (jean-francois@famillegratton.net)
- Enhanced server version (jean-francois@famillegratton.net)
- completed srv reload and srv version (jean-francois@famillegratton.net)
- minor types refactoring (jean-francois@famillegratton.net)

* Thu Sep 11 2025 Binary package builder <builder@famillegratton.net> 1.30.00-0
- version bump (jean-francois@famillegratton.net)
- ide config update (jean-francois@famillegratton.net)
- completed db subcommands (jean-francois@famillegratton.net)
- module name refactoring (jean-francois@famillegratton.net)
- types subpackage refactoring (jean-francois@famillegratton.net)
- Revamped environment handling (jean-francois@famillegratton.net)

* Tue Sep 09 2025 Binary package builder <builder@famillegratton.net> 1.21.10-0
- new package built with tito

* Sun Sep 07 2025 Binary package builder <builder@famillegratton.net> 1.21.10-0
- Updated to GO 1.25.1 (jean-francois@famillegratton.net)

* Thu Jul 24 2025 Binary package builder <builder@famillegratton.net> 1.21.00-0
- Fixed nil-pointer issue in db backup -u (jean-francois@famillegratton.net)

* Thu Jul 24 2025 Binary package builder <builder@famillegratton.net> 1.20.00-0
- Software version bump (jean-francois@famillegratton.net)
- Phase 1 of making error codes more consistent (jean-
  francois@famillegratton.net)
- migrated customError to v2 (jean-francois@famillegratton.net)
- Moved the logger in its own package (jean-francois@famillegratton.net)

* Tue Jul 15 2025 Binary package builder <builder@famillegratton.net> 1.10.00-0
- updated gitignore (jean-francois@famillegratton.net)
- removed samples (jean-francois@famillegratton.net)
- Tool is now complete (jean-francois@famillegratton.net)
- Version bump (jean-francois@famillegratton.net)
- Fixed double-quote mixup, added constraints management (jean-
  francois@famillegratton.net)

* Mon Jul 14 2025 Binary package builder <builder@famillegratton.net> 1.07.00-1
- Release bump (jean-francois@famillegratton.net)

* Mon Jul 14 2025 Binary package builder <builder@famillegratton.net> 1.07.00-0
- Version bump (jean-francois@famillegratton.net)
- Added PK management (jean-francois@famillegratton.net)
- fixed wrong command sequence when restoring the db (jean-
  francois@famillegratton.net)
- fixed issue with sslmode not returning the correct value in connection string
  (jean-francois@famillegratton.net)
- updated configs and samples (jean-francois@famillegratton.net)
- Fixed restore failing because of \c; added ownership and attributes to backup
  (jean-francois@famillegratton.net)
- interim sync (jean-francois@famillegratton.net)
- Version bump (jean-francois@famillegratton.net)
- fixed script filemode (builder@famillegratton.net)

* Wed Jul 09 2025 Binary package builder <builder@famillegratton.net> 1.00.00-0
- new package built with tito

